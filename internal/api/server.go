// Package api is the HTTP surface: JSON endpoints under /api plus the
// static single-page UI. Handlers are thin: validate, delegate to a store,
// respond. Stores are interfaces so tests inject fakes.
package api

import (
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type Server struct {
	notes NoteStore  // nil when no database is bound
	blobs BlobStore  // nil when no bucket is bound
	cache CacheStore // nil when no Valkey is bound
	uiDir string
}

func New(notes NoteStore, blobs BlobStore, cache CacheStore, uiDir string) http.Handler {
	s := &Server{notes: notes, blobs: blobs, cache: cache, uiDir: uiDir}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", s.health)
	mux.HandleFunc("GET /api/visits", s.visits)
	mux.HandleFunc("GET /api/notes", s.listNotes)
	mux.HandleFunc("POST /api/notes", s.createNote)
	mux.HandleFunc("DELETE /api/notes/{id}", s.deleteNote)
	mux.HandleFunc("GET /api/files", s.listFiles)
	mux.HandleFunc("POST /api/files", s.uploadFile)
	mux.HandleFunc("GET /api/files/{key}", s.downloadFile)
	mux.HandleFunc("DELETE /api/files/{key}", s.deleteFile)
	mux.Handle("/", s.spa())

	return logRequests(mux)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":   "ok",
		"database": s.notes != nil,
		"storage":  s.blobs != nil,
		"cache":    s.cache != nil,
	})
}

// spa serves the built SvelteKit app from uiDir, falling back to index.html
// for client-side routes. When the UI is not built (API-only dev), it
// returns a hint instead of a blank 404.
func (s *Server) spa() http.Handler {
	if _, err := os.Stat(s.uiDir); err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writeError(w, http.StatusNotFound,
				"UI not built - run `make dev-ui` for the dev server or `make build` to build it")
		})
	}
	fsys := os.DirFS(s.uiDir)
	files := http.FileServerFS(fsys)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name != "" {
			if _, err := fs.Stat(fsys, name); err == nil {
				files.ServeHTTP(w, r)
				return
			}
		}
		http.ServeFileFS(w, r, fsys, "index.html")
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(rec, r)
		if strings.HasPrefix(r.URL.Path, "/api/") {
			slog.Info("request", "method", r.Method, "path", r.URL.Path,
				"status", rec.status, "duration", time.Since(start).Round(time.Millisecond))
		}
	})
}
