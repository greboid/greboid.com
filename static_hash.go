package main

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

type StaticHashManager struct {
	hashes map[string]string
	files  map[string][]byte
}

func InitStaticHashManager(embedFS embed.FS, staticDir string) (*StaticHashManager, error) {
	mgr := &StaticHashManager{
		hashes: make(map[string]string),
		files:  make(map[string][]byte),
	}

	err := fs.WalkDir(embedFS, staticDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		data, err := fs.ReadFile(embedFS, path)
		if err != nil {
			return fmt.Errorf("reading file %s: %w", path, err)
		}

		hash := computeFileHash(data)

		relPath, err := filepath.Rel(staticDir, path)
		if err != nil {
			return fmt.Errorf("getting relative path for %s: %w", path, err)
		}

		hashedFilename := generateHashedFilename(relPath, hash)

		mgr.hashes[relPath] = hashedFilename
		mgr.files[hashedFilename] = data

		slog.Debug("loaded static file", "original", relPath, "hashed", hashedFilename)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking static directory: %w", err)
	}

	if len(mgr.hashes) == 0 {
		slog.Warn("no static files found", "directory", staticDir)
	}

	return mgr, nil
}

func (s *StaticHashManager) GetHashedFilename(original string) string {
	return s.hashes[original]
}

func (s *StaticHashManager) ServeHashedFile(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/static/")
	if filename == "" {
		http.NotFound(w, r)
		return
	}

	content, exists := s.files[filename]
	if !exists {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(filename)))

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(content); err != nil {
		slog.Error("error writing response", "filename", filename, "error", err)
	}
}

func computeFileHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])[:8]
}

func generateHashedFilename(original, hash string) string {
	ext := filepath.Ext(original)
	base := strings.TrimSuffix(original, ext)
	return fmt.Sprintf("%s-%s%s", base, hash, ext)
}
