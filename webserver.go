package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/csmith/middleware"
)

type TemplateData struct {
	Title        string
	StaticFiles  map[string]string
	PreviousLink string
	Template     string
}

type WebServer struct {
	webFS         embed.FS
	staticHashMgr *StaticHashManager
	staticFiles   map[string]string
	mux           *http.ServeMux
	chain         http.Handler
}

func InitWebServer(embedFS embed.FS, shm *StaticHashManager) (*WebServer, error) {
	ws := &WebServer{
		webFS:         embedFS,
		staticHashMgr: shm,
		staticFiles:   make(map[string]string),
		mux:           http.NewServeMux(),
	}
	ws.chain = middleware.Chain(middleware.WithMiddleware(
		middleware.CacheControl(),
		middleware.ErrorHandler(
			middleware.WithErrorHandler(http.StatusNotFound, http.HandlerFunc(ws.serve404)),
			middleware.WithErrorHandler(http.StatusInternalServerError, http.HandlerFunc(ws.serve500)),
		),
		middleware.RealAddress(),
		middleware.Compress(),
		middleware.Recover(middleware.WithPanicLogger(func(r *http.Request, err any) {
			slog.Error("Panic recovered from: %w", err)
		})),
	))(ws.mux)
	ws.addStaticFiles()
	ws.addRoutes()
	return ws, nil
}

func (ws *WebServer) ListenAndServe(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), ws.chain)
}

func (ws *WebServer) addStaticFiles() {
	ws.staticFiles["MainCSS"] = ws.staticHashMgr.GetHashedFilename("main.css")
	ws.staticFiles["Favicon"] = ws.staticHashMgr.GetHashedFilename("favicon.svg")
	ws.staticFiles["TransitionsJS"] = ws.staticHashMgr.GetHashedFilename("transitions.js")
}

func (ws *WebServer) addRoutes() {
	ws.mux.HandleFunc("/favicon.ico", ws.faviconHandler)
	ws.mux.HandleFunc("/me/", ws.meHandler)
	ws.mux.HandleFunc("/static/", ws.staticHashMgr.ServeHashedFile)
	ws.mux.HandleFunc("/", ws.rootHandler)
}

func (ws *WebServer) renderTemplate(w http.ResponseWriter, templatePath string, data TemplateData, status int) error {
	tmpl, err := template.ParseFS(ws.webFS, "web/layouts/base.html", templatePath)
	if err != nil {
		return fmt.Errorf("parsing template %s: %w", templatePath, err)
	}

	data.StaticFiles = ws.staticFiles
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.WriteHeader(status)

	if err = tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		return fmt.Errorf("executing template %s: %w", templatePath, err)
	}

	return nil
}

func (ws *WebServer) serve404(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title: "Page Not Found",
	}
	slog.Debug("Page not found", "page", r.URL)
	if err := ws.renderTemplate(w, "web/pages/404.html", data, http.StatusNotFound); err != nil {
		slog.Error("error rendering 404 template", "error", err)
		http.NotFound(w, r)
	}
}

func (ws *WebServer) serve500(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title: "Error serving page",
	}

	if err := ws.renderTemplate(w, "web/pages/500.html", data, http.StatusInternalServerError); err != nil {
		slog.Error("error rendering 500 template", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (ws *WebServer) faviconHandler(w http.ResponseWriter, r *http.Request) {
	data, err := fs.ReadFile(ws.webFS, "web/static/favicon.ico")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/x-icon")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		slog.Error("error writing favicon", "error", err)
	}
}

func (ws *WebServer) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		data := TemplateData{
			Title: "Greg Holmes",
		}
		if err := ws.renderTemplate(w, "web/pages/index.html", data, http.StatusOK); err != nil {
			slog.Error("error rendering index template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	path := filepath.Join("web", r.URL.Path)
	info, err := fs.Stat(ws.webFS, path)
	if err == nil && !info.IsDir() {
		http.ServeFileFS(w, r, ws.webFS, path)
		return
	}

	http.NotFound(w, r)
}

func (ws *WebServer) meHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/me" || r.URL.Path == "/me/" {
		http.NotFound(w, r)
		return
	}
	templateName := strings.TrimPrefix(r.URL.Path, "/me/")
	templateName = strings.TrimSuffix(templateName, ".html")
	templatePath := filepath.Join("web/pages/me", templateName+".html")

	data := TemplateData{
		Title:        "Greg Holmes",
		PreviousLink: "/",
	}
	if err := ws.renderTemplate(w, templatePath, data, http.StatusOK); err == nil {
		return
	}
	path := filepath.Join("web", r.URL.Path)

	if info, err := fs.Stat(ws.webFS, path); err == nil && info.IsDir() {
		http.NotFound(w, r)
		return
	}

	if _, err := fs.Stat(ws.webFS, path); err != nil {
		htmlPath := path + ".html"
		if info, err := fs.Stat(ws.webFS, htmlPath); err == nil && !info.IsDir() {
			http.ServeFileFS(w, r, ws.webFS, htmlPath)
			return
		}
		http.NotFound(w, r)
		return
	}

	http.ServeFileFS(w, r, ws.webFS, path)
}
