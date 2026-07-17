package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// fakeCache is an in-memory CacheStore for handler tests.
type fakeCache struct {
	counters map[string]int64
	values   map[string][]byte
}

func newFakeCache() *fakeCache {
	return &fakeCache{counters: map[string]int64{}, values: map[string][]byte{}}
}

func (f *fakeCache) Incr(ctx context.Context, key string) (int64, error) {
	f.counters[key]++
	return f.counters[key], nil
}

func (f *fakeCache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	v, ok := f.values[key]
	return v, ok, nil
}

func (f *fakeCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	f.values[key] = value
	return nil
}

func (f *fakeCache) Del(ctx context.Context, key string) error {
	delete(f.values, key)
	return nil
}

func TestVisitsCounter(t *testing.T) {
	srv := New(nil, nil, newFakeCache(), "does-not-exist")

	for want := 1; want <= 3; want++ {
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, httptest.NewRequest("GET", "/api/visits", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("visits: status = %d, want 200 (%s)", rec.Code, rec.Body)
		}
		if got := rec.Body.String(); !strings.Contains(got, `"visits":`+string(rune('0'+want))) {
			t.Fatalf("visits = %s, want %d", got, want)
		}
	}
}

func TestVisitsWithoutCacheReturnsHint(t *testing.T) {
	srv := New(nil, nil, nil, "does-not-exist")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest("GET", "/api/visits", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rec.Code)
	}
}

func TestNotesListUsesCacheAside(t *testing.T) {
	store := &fakeNotes{}
	c := newFakeCache()
	srv := New(store, nil, c, "does-not-exist")

	// First list: miss (fills the cache).
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest("GET", "/api/notes", nil))
	if got := rec.Header().Get("X-Cache"); got != "miss" {
		t.Fatalf("first list X-Cache = %q, want miss", got)
	}

	// Second list: hit.
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest("GET", "/api/notes", nil))
	if got := rec.Header().Get("X-Cache"); got != "hit" {
		t.Fatalf("second list X-Cache = %q, want hit", got)
	}

	// A write invalidates: next list is a miss again.
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest("POST", "/api/notes", strings.NewReader(`{"text":"x"}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: status = %d (%s)", rec.Code, rec.Body)
	}
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest("GET", "/api/notes", nil))
	if got := rec.Header().Get("X-Cache"); got != "miss" {
		t.Fatalf("post-write list X-Cache = %q, want miss", got)
	}
}
