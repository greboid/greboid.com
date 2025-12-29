package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/csmith/middleware"
)

type TemplateData struct {
	Title        string
	StaticFiles  map[string]string
	PreviousLink string
	Template     string
}

type WebServer struct {
	webFS         fs.FS
	staticHashMgr *StaticHashManager
	staticFiles   map[string]string
	mux           *http.ServeMux
	chain         http.Handler
}

func NewWebServer(webFS fs.FS, shm *StaticHashManager) (*WebServer, error) {
	ws := &WebServer{
		webFS:         webFS,
		staticHashMgr: shm,
		staticFiles:   make(map[string]string),
		mux:           http.NewServeMux(),
	}
	ws.chain = middleware.Chain(middleware.WithMiddleware(
		middleware.StripTrailingSlashes(),
		middleware.CacheControl(),
		middleware.ErrorHandler(
			middleware.WithErrorHandler(http.StatusNotFound, http.HandlerFunc(ws.serve404)),
			middleware.WithErrorHandler(http.StatusInternalServerError, http.HandlerFunc(ws.serve500)),
		),
		middleware.RealAddress(),
		middleware.Compress(),
		middleware.Recover(middleware.WithPanicLogger(func(r *http.Request, err any) {
			slog.Error("panic recovered", "error", err, "path", r.URL.Path, "method", r.Method)
		})),
	))(ws.mux)
	ws.addStaticFiles()
	ws.addRoutes()
	return ws, nil
}

func (ws *WebServer) ListenAndServe(port int) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: ws.chain,
	}
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("starting server", "port", port)
		serverErrors <- server.ListenAndServe()
	}()
	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		slog.Info("shutdown signal received", "signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			_ = server.Close()
			return fmt.Errorf("could not gracefully shutdown server: %w", err)
		}

		slog.Info("server shutdown complete")
	}

	return nil
}

func (ws *WebServer) addStaticFiles() {
	ws.staticFiles["MainCSS"] = ws.staticHashMgr.GetHashedFilename("main.css")
	ws.staticFiles["Favicon"] = ws.staticHashMgr.GetHashedFilename("favicon.svg")
	ws.staticFiles["Dots"] = ws.staticHashMgr.GetHashedFilename("dots.svg")
}

func (ws *WebServer) addRoutes() {
	ws.mux.HandleFunc("/favicon.ico", ws.faviconHandler)
	ws.mux.HandleFunc("/me/{name}", ws.meHandler)
	ws.mux.HandleFunc("/static/"+ws.staticFiles["MainCSS"], ws.cssHandler)
	ws.mux.HandleFunc("/static/", ws.staticHashMgr.ServeHashedFile)
	ws.mux.HandleFunc("/{$}", ws.rootHandler)
}

func (ws *WebServer) renderTemplateOrError(w http.ResponseWriter, _ *http.Request, templatePath string, data TemplateData, statusCode int) {
	_, err := fs.Stat(ws.webFS, templatePath)
	if errors.Is(err, fs.ErrNotExist) {
		http.NotFound(w, nil)
		return
	}

	tmpl, err := template.ParseFS(ws.webFS, "layouts/base.html", templatePath)
	if err != nil {
		logAndError(w, "parsing template", http.StatusInternalServerError, err, "templatePath", templatePath)
		return
	}

	data.StaticFiles = ws.staticFiles

	var buf bytes.Buffer
	if err = tmpl.ExecuteTemplate(&buf, "base.html", data); err != nil {
		logAndError(w, "executing template", http.StatusInternalServerError, err, "templatePath", templatePath)
		return
	}

	writeResponse(w, "text/html; charset=utf-8", statusCode, buf.Bytes())
}

func writeResponse(w http.ResponseWriter, contentType string, statusCode int, data []byte) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	if _, err := w.Write(data); err != nil {
		slog.Error("error writing response", "error", err)
	}
}

func logAndError(w http.ResponseWriter, message string, status int, err error, logAttrs ...any) {
	slog.Error(message, append([]any{"error", err}, logAttrs...)...)
	http.Error(w, http.StatusText(status), status)
}

func (ws *WebServer) readFileOrNotFound(w http.ResponseWriter, r *http.Request, path string) ([]byte, bool) {
	data, err := fs.ReadFile(ws.webFS, path)
	if err != nil {
		http.NotFound(w, r)
		return nil, false
	}
	return data, true
}

func (ws *WebServer) serve404(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Page not found", "page", r.URL)
	ws.renderTemplateOrError(w, r, "pages/404.html", TemplateData{
		Title: "Page Not Found",
	}, http.StatusNotFound)
}

func (ws *WebServer) serve500(w http.ResponseWriter, r *http.Request) {
	ws.renderTemplateOrError(w, r, "pages/500.html", TemplateData{
		Title: "Error serving page",
	}, http.StatusInternalServerError)
}

func (ws *WebServer) cssHandler(w http.ResponseWriter, r *http.Request) {
	cssData, ok := ws.readFileOrNotFound(w, r, "static/main.css")
	if !ok {
		return
	}

	tmpl, err := template.New("css").Parse(string(cssData))
	if err != nil {
		logAndError(w, "error parsing CSS template", http.StatusInternalServerError, err)
		return
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, ws.staticFiles); err != nil {
		logAndError(w, "error executing CSS template", http.StatusInternalServerError, err)
		return
	}

	writeResponse(w, "text/css; charset=utf-8", http.StatusOK, buf.Bytes())
}

func (ws *WebServer) faviconHandler(w http.ResponseWriter, r *http.Request) {
	data, ok := ws.readFileOrNotFound(w, r, "static/favicon.ico")
	if !ok {
		return
	}
	writeResponse(w, "image/x-icon", http.StatusOK, data)
}

func (ws *WebServer) rootHandler(w http.ResponseWriter, r *http.Request) {
	ws.renderTemplateOrError(w, r, "pages/index.html", TemplateData{
		Title: "Greg Holmes",
	}, http.StatusOK)
}

func (ws *WebServer) meHandler(w http.ResponseWriter, r *http.Request) {
	templateName := r.PathValue("name")
	if strings.HasSuffix(templateName, ".html") {
		templateName = strings.TrimSuffix(templateName, ".html")
	}
	templatePath := filepath.Join("pages/me", templateName+".html")

	ws.renderTemplateOrError(w, r, templatePath, TemplateData{
		Title:        "Greg Holmes: " + cases.Title(language.BritishEnglish).String(templateName),
		PreviousLink: "/",
	}, http.StatusOK)
}
