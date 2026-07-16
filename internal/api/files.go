package api

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"time"
)

type FileInfo struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
}

// ErrNotFound is returned by BlobStore.GetFile when the key does not exist.
var ErrNotFound = errors.New("not found")

type BlobStore interface {
	ListFiles(ctx context.Context) ([]FileInfo, error)
	PutFile(ctx context.Context, key, contentType string, r io.Reader, size int64) error
	GetFile(ctx context.Context, key string) (io.ReadCloser, string, error)
	DeleteFile(ctx context.Context, key string) error
}

const maxUploadBytes = 25 << 20 // 25 MiB

const noStorageHint = "no storage bound - add a Bucket service to your klickops project and accept the suggested binding (locally: `docker compose up -d` and copy .env.example to .env)"

func (s *Server) listFiles(w http.ResponseWriter, r *http.Request) {
	if s.blobs == nil {
		writeError(w, http.StatusServiceUnavailable, noStorageHint)
		return
	}
	files, err := s.blobs.ListFiles(r.Context())
	if err != nil {
		internalError(w, "failed to list files", err)
		return
	}
	if files == nil {
		files = []FileInfo{}
	}
	writeJSON(w, http.StatusOK, files)
}

func (s *Server) uploadFile(w http.ResponseWriter, r *http.Request) {
	if s.blobs == nil {
		writeError(w, http.StatusServiceUnavailable, noStorageHint)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "send a multipart form with a \"file\" field (max 25 MiB)")
		return
	}
	defer file.Close()

	key := path.Base(header.Filename)
	if key == "" || key == "." || key == "/" {
		writeError(w, http.StatusBadRequest, "file must have a name")
		return
	}
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	if err := s.blobs.PutFile(r.Context(), key, contentType, file, header.Size); err != nil {
		internalError(w, "failed to store file", err)
		return
	}
	writeJSON(w, http.StatusCreated, FileInfo{Key: key, Size: header.Size, LastModified: time.Now().UTC()})
}

func (s *Server) downloadFile(w http.ResponseWriter, r *http.Request) {
	if s.blobs == nil {
		writeError(w, http.StatusServiceUnavailable, noStorageHint)
		return
	}
	key, ok := fileKey(w, r)
	if !ok {
		return
	}
	body, contentType, err := s.blobs.GetFile(r.Context(), key)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "file not found")
		return
	}
	if err != nil {
		internalError(w, "failed to fetch file", err)
		return
	}
	defer body.Close()
	w.Header().Set("Content-Type", contentType)
	if _, err := io.Copy(w, body); err != nil {
		slog.Error("failed to stream file", "key", key, "error", err)
	}
}

func (s *Server) deleteFile(w http.ResponseWriter, r *http.Request) {
	if s.blobs == nil {
		writeError(w, http.StatusServiceUnavailable, noStorageHint)
		return
	}
	key, ok := fileKey(w, r)
	if !ok {
		return
	}
	if err := s.blobs.DeleteFile(r.Context(), key); err != nil {
		internalError(w, "failed to delete file", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func fileKey(w http.ResponseWriter, r *http.Request) (string, bool) {
	key := r.PathValue("key")
	if key == "" || strings.ContainsAny(key, "/\\") || strings.Contains(key, "..") {
		writeError(w, http.StatusBadRequest, "invalid file name")
		return "", false
	}
	return key, true
}

func internalError(w http.ResponseWriter, message string, err error) {
	slog.Error(message, "error", err)
	writeError(w, http.StatusInternalServerError, message)
}
