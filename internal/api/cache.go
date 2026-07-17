package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// CacheStore is what the handlers need from a cache. Implemented by
// internal/cache (Valkey); tests inject fakes.
type CacheStore interface {
	Incr(ctx context.Context, key string) (int64, error)
	Get(ctx context.Context, key string) ([]byte, bool, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

const noCacheHint = "no cache bound - add a Valkey service to your klickops project and connect it to this app (it injects REDIS_URL; locally: `docker compose up -d` and copy .env.example to .env)"

const (
	visitsKey     = "starter:visits"
	notesCacheKey = "starter:notes"
	notesCacheTTL = 30 * time.Second
)

// visits counts page loads with an atomic INCR — the classic first
// Valkey use case.
func (s *Server) visits(w http.ResponseWriter, r *http.Request) {
	if s.cache == nil {
		writeError(w, http.StatusServiceUnavailable, noCacheHint)
		return
	}
	n, err := s.cache.Incr(r.Context(), visitsKey)
	if err != nil {
		internalError(w, "failed to count visit", err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"visits": n})
}

// cachedNotes is the cache-aside pattern: serve the notes list from
// Valkey when fresh, fall through to PostgreSQL on a miss, and let
// writes invalidate (see createNote / deleteNote). The X-Cache header
// makes hits observable in the browser dev tools.
func (s *Server) cachedNotes(ctx context.Context) ([]Note, bool, error) {
	if s.cache == nil {
		notes, err := s.notes.ListNotes(ctx)
		return notes, false, err
	}
	if raw, ok, err := s.cache.Get(ctx, notesCacheKey); err == nil && ok {
		var notes []Note
		if json.Unmarshal(raw, &notes) == nil {
			return notes, true, nil
		}
	}
	notes, err := s.notes.ListNotes(ctx)
	if err != nil {
		return nil, false, err
	}
	if raw, err := json.Marshal(notes); err == nil {
		// Best-effort: a failed cache write must never fail the request.
		_ = s.cache.Set(ctx, notesCacheKey, raw, notesCacheTTL)
	}
	return notes, false, nil
}

// invalidateNotes drops the cached list after a write. Best-effort — at
// worst the stale list survives until the TTL expires.
func (s *Server) invalidateNotes(ctx context.Context) {
	if s.cache != nil {
		_ = s.cache.Del(ctx, notesCacheKey)
	}
}
