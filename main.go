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

var (
	port    = flag.Int("port", 8080, "Port for webserver to listen on")
	devMode = flag.Bool("dev", false, "Enable development mode with live reload")
)

func getWebFS() (fs.FS, error) {
	if _, err := os.Stat("web"); err == nil {
		slog.Info("using web directory from disk")
		return os.DirFS("web"), nil
	}
	slog.Info("using embedded web directory")
	return fs.Sub(webFS, "web")
}

func main() {
	envflag.Parse()
	slogflags.Logger(slogflags.WithSetDefault(true))
	slog.Debug("Debug is enabled")

	subFS, err := getWebFS()
	if err != nil {
		slog.Error("failed to get web filesystem", "error", err)
		os.Exit(1)
	}

	staticHashMgr, err := NewStaticHashManager(subFS, "static")
	if err != nil {
		slog.Error("failed to initialize static hash manager", "error", err)
		os.Exit(1)
	}

	var reloadMgr *ReloadManager
	if *devMode {
		slog.Info("development mode enabled - starting file watcher")
		reloadMgr, err = NewReloadManager("web")
		if err != nil {
			slog.Error("failed to initialize reload manager", "error", err)
			os.Exit(1)
		}
		defer func(reloadMgr *ReloadManager) {
			_ = reloadMgr.Close()
		}(reloadMgr)
	}

	ws, err := NewWebServer(subFS, staticHashMgr, *devMode, reloadMgr)
	if err != nil {
		slog.Error("failed to initialize web server", "error", err)
		os.Exit(1)
	}

	slog.Info("starting server", "address", fmt.Sprintf("http://localhost:%d", *port))
	if err = ws.ListenAndServe(*port); err != nil {
		slog.Error("server error", "error", err)
	}
}
