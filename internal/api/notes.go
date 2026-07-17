package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Note struct {
	ID        int64     `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

type NoteStore interface {
	ListNotes(ctx context.Context) ([]Note, error)
	CreateNote(ctx context.Context, text string) (Note, error)
	DeleteNote(ctx context.Context, id int64) error
}

const maxNoteLength = 2000

const noDatabaseHint = "no database bound - add a PostgreSQL service to your klickops project and bind it as DATABASE_URL (locally: `docker compose up -d` and copy .env.example to .env)"

func (s *Server) listNotes(w http.ResponseWriter, r *http.Request) {
	if s.notes == nil {
		writeError(w, http.StatusServiceUnavailable, noDatabaseHint)
		return
	}
	notes, hit, err := s.cachedNotes(r.Context())
	if err != nil {
		internalError(w, "failed to list notes", err)
		return
	}
	if notes == nil {
		notes = []Note{}
	}
	if s.cache != nil {
		if hit {
			w.Header().Set("X-Cache", "hit")
		} else {
			w.Header().Set("X-Cache", "miss")
		}
	}
	writeJSON(w, http.StatusOK, notes)
}

func (s *Server) createNote(w http.ResponseWriter, r *http.Request) {
	if s.notes == nil {
		writeError(w, http.StatusServiceUnavailable, noDatabaseHint)
		return
	}
	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "request body must be JSON like {\"text\": \"...\"}")
		return
	}
	req.Text = strings.TrimSpace(req.Text)
	if req.Text == "" {
		writeError(w, http.StatusBadRequest, "text is required")
		return
	}
	if utf8.RuneCountInString(req.Text) > maxNoteLength {
		writeError(w, http.StatusBadRequest, "text must be at most 2000 characters")
		return
	}
	note, err := s.notes.CreateNote(r.Context(), req.Text)
	if err != nil {
		internalError(w, "failed to create note", err)
		return
	}
	s.invalidateNotes(r.Context())
	writeJSON(w, http.StatusCreated, note)
}

func (s *Server) deleteNote(w http.ResponseWriter, r *http.Request) {
	if s.notes == nil {
		writeError(w, http.StatusServiceUnavailable, noDatabaseHint)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id must be a number")
		return
	}
	if err := s.notes.DeleteNote(r.Context(), id); err != nil {
		internalError(w, "failed to delete note", err)
		return
	}
	s.invalidateNotes(r.Context())
	w.WriteHeader(http.StatusNoContent)
}
