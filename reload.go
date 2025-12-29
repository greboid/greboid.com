package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const DebounceTimeout = 100 * time.Millisecond

type ReloadManager struct {
	watcher   *fsnotify.Watcher
	clients   map[chan struct{}]bool
	mu        sync.RWMutex
	debouncer *time.Timer
}

func NewReloadManager(webDir string) (*ReloadManager, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	rm := &ReloadManager{
		watcher: watcher,
		clients: make(map[chan struct{}]bool),
	}

	if err = rm.watchRecursive(webDir); err != nil {
		_ = watcher.Close()
		return nil, err
	}

	go rm.watch()

	return rm, nil
}

func (rm *ReloadManager) watchRecursive(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if err = rm.watcher.Add(path); err != nil {
			slog.Warn("failed to watch path", "path", path, "error", err)
		}
		return nil
	})
}

func (rm *ReloadManager) watch() {
	for {
		select {
		case event, ok := <-rm.watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) ||
				event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				rm.debounceReload()
			}
		case err, ok := <-rm.watcher.Errors:
			if !ok {
				return
			}
			slog.Error("watcher error", "error", err)
		}
	}
}

func (rm *ReloadManager) debounceReload() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if rm.debouncer != nil {
		rm.debouncer.Stop()
	}

	rm.debouncer = time.AfterFunc(DebounceTimeout, func() {
		rm.notifyClients()
	})
}

func (rm *ReloadManager) notifyClients() {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	for client := range rm.clients {
		select {
		case client <- struct{}{}:
		default:
		}
	}
}

func (rm *ReloadManager) Subscribe() chan struct{} {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	client := make(chan struct{}, 1)
	rm.clients[client] = true
	return client
}

func (rm *ReloadManager) Unsubscribe(client chan struct{}) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	delete(rm.clients, client)
	close(client)
}

func (rm *ReloadManager) Close() error {
	return rm.watcher.Close()
}

func (rm *ReloadManager) ServeSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	client := rm.Subscribe()
	defer rm.Unsubscribe(client)

	_, _ = fmt.Fprintf(w, "data: connected\n\n")
	flusher.Flush()
	for {
		select {
		case <-client:
			_, _ = fmt.Fprintf(w, "data: reload\n\n")
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
