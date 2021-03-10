package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//go:embed static
var staticfs embed.FS

func main() {
	staticFiles, err := fs.Sub(staticfs, "static")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	router := mux.NewRouter()
	router.Use(handlers.ProxyHeaders)
	router.Use(handlers.CompressHandler)
	router.Use(NewLoggingHandler(os.Stdout))
	router.PathPrefix("/").Handler(NotFoundHandler(CheckWebP(http.FileServer(http.FS(staticFiles)), staticFiles), staticFiles))

	log.Print("Starting server.")
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		_ = server.ListenAndServe()
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Unable to shutdown: %s", err.Error())
	}
	log.Print("Finishing server.")
}

func CheckWebP(h http.Handler, files fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept"), "image/webp") {
			webp := r.URL.Path + ".webp"
			_, err := files.Open(strings.TrimPrefix(webp, "/"))
			if err == nil {
				r.URL.Path = webp
			}
		}
		h.ServeHTTP(w, r)
	}
}