package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeNotes struct {
	notes  []Note
	nextID int64
}

func (f *fakeNotes) ListNotes(ctx context.Context) ([]Note, error) { return f.notes, nil }

func (f *fakeNotes) CreateNote(ctx context.Context, text string) (Note, error) {
	f.nextID++
	n := Note{ID: f.nextID, Text: text}
	f.notes = append(f.notes, n)
	return n, nil
}

func (f *fakeNotes) DeleteNote(ctx context.Context, id int64) error {
	for i, n := range f.notes {
		if n.ID == id {
			f.notes = append(f.notes[:i], f.notes[i+1:]...)
			return nil
		}
	}
	return nil
}

func newTestServer(notes NoteStore) http.Handler {
	return New(notes, nil, nil, "does-not-exist")
}

func TestNotesCRUD(t *testing.T) {
	store := &fakeNotes{}
	srv := newTestServer(store)

	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest("POST", "/api/notes", strings.NewReader(`{"text":"hello"}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: status = %d, want 201 (%s)", rec.Code, rec.Body)
	}
	var created Note
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("create: bad body: %v", err)
	}
	if created.Text != "hello" {
		t.Errorf("create: text = %q", created.Text)
	}

	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest("GET", "/api/notes", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("list: status = %d", rec.Code)
	}
	var listed []Note
	if err := json.NewDecoder(rec.Body).Decode(&listed); err != nil {
		t.Fatalf("list: bad body: %v", err)
	}
	if len(listed) != 1 {
		t.Fatalf("list: got %d notes, want 1", len(listed))
	}

	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest("DELETE", "/api/notes/1", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete: status = %d, want 204", rec.Code)
	}
	if len(store.notes) != 0 {
		t.Errorf("delete: store still has %d notes", len(store.notes))
	}
}

func TestCreateNoteValidation(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"empty text", `{"text":""}`},
		{"whitespace only", `{"text":"   "}`},
		{"not json", `hello`},
		{"too long", `{"text":"` + strings.Repeat("x", 2001) + `"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			newTestServer(&fakeNotes{}).ServeHTTP(rec,
				httptest.NewRequest("POST", "/api/notes", strings.NewReader(tt.body)))
			if rec.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", rec.Code)
			}
		})
	}
}

func TestNotesWithoutDatabase(t *testing.T) {
	srv := newTestServer(nil)
	for _, req := range []*http.Request{
		httptest.NewRequest("GET", "/api/notes", nil),
		httptest.NewRequest("POST", "/api/notes", strings.NewReader(`{"text":"x"}`)),
		httptest.NewRequest("DELETE", "/api/notes/1", nil),
	} {
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("%s %s: status = %d, want 503", req.Method, req.URL.Path, rec.Code)
		}
	}
}

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestServer(&fakeNotes{}).ServeHTTP(rec, httptest.NewRequest("GET", "/api/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var health struct {
		Status   string `json:"status"`
		Database bool   `json:"database"`
		Storage  bool   `json:"storage"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&health); err != nil {
		t.Fatal(err)
	}
	if health.Status != "ok" || !health.Database || health.Storage {
		t.Errorf("health = %+v", health)
	}
}
