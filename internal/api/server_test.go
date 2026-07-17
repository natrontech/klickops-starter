package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSPAFallback(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html>app</html>"), 0o644)
	os.WriteFile(filepath.Join(dir, "app.css"), []byte("body{}"), 0o644)

	srv := New(nil, nil, nil, dir)

	tests := []struct {
		path     string
		wantBody string
	}{
		{"/", "<html>app</html>"},
		{"/app.css", "body{}"},
		{"/some/client/route", "<html>app</html>"},
	}
	for _, tt := range tests {
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, httptest.NewRequest("GET", tt.path, nil))
		if rec.Code != http.StatusOK {
			t.Errorf("%s: status = %d", tt.path, rec.Code)
		}
		if !strings.Contains(rec.Body.String(), tt.wantBody) {
			t.Errorf("%s: body = %q, want %q", tt.path, rec.Body.String(), tt.wantBody)
		}
	}
}

func TestSPAMissingUIDir(t *testing.T) {
	rec := httptest.NewRecorder()
	New(nil, nil, nil, "does-not-exist").ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "UI not built") {
		t.Errorf("body = %q, want UI-not-built hint", rec.Body.String())
	}
}
