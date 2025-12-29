package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"

	"github.com/csmith/envflag/v2"
	"github.com/csmith/slogflags"
)

//go:embed web
var webFS embed.FS

var port = flag.Int("port", 8080, "Port for webserver to listen on")

func main() {
	envflag.Parse()
	slogflags.Logger(slogflags.WithSetDefault(true))
	slog.Debug("Debug is enabled")

	subFS, err := fs.Sub(webFS, "web")
	if err != nil {
		slog.Error("failed to create sub-filesystem", "error", err)
		os.Exit(1)
	}

	staticHashMgr, err := NewStaticHashManager(subFS, "static")
	if err != nil {
		slog.Error("failed to initialize static hash manager", "error", err)
		os.Exit(1)
	}
	ws, err := NewWebServer(subFS, staticHashMgr)
	if err != nil {
		slog.Error("failed to initialize web server", "error", err)
		os.Exit(1)
	}

	slog.Info("starting server", "address", fmt.Sprintf("http://localhost:%d", *port))
	if err = ws.ListenAndServe(*port); err != nil {
		slog.Error("server error", "error", err)
	}
}
