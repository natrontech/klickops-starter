---
name: new-feature
description: Add a full-stack feature (API endpoint + database table + UI) to this app. Use whenever the user asks for a new capability that stores or serves data - e.g. "add tasks", "let users upload avatars", "add a contact form".
---

# Adding a full-stack feature

Follow the notes feature as the reference implementation end-to-end:
migration `internal/db/migrations/0001_notes.sql` → store
`internal/db/notes.go` → handlers `internal/api/notes.go` → tests
`internal/api/notes_test.go` → frontend `ui/src/lib/api/notes.ts` →
UI `ui/src/routes/+page.svelte`.

Steps for a new resource called `things`:

## 1. Migration

Create `internal/db/migrations/0002_things.sql` (next free number, never
edit an applied migration):

```sql
CREATE TABLE IF NOT EXISTS things (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

It is applied automatically at startup.

## 2. Types + store interface (internal/api/things.go)

Define the wire type with camelCase JSON tags and the interface the
handlers consume:

```go
type Thing struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"createdAt"`
}

type ThingStore interface {
    ListThings(ctx context.Context) ([]Thing, error)
    CreateThing(ctx context.Context, name string) (Thing, error)
    DeleteThing(ctx context.Context, id int64) error
}
```

## 3. Store implementation (internal/db/things.go)

Plain SQL through pgx, parameterized queries only, errors wrapped with
`%w`. Mirror `internal/db/notes.go`.

## 4. Handlers + routes

Handlers in `internal/api/things.go`: nil-store guard (503 + hint),
validate input (400 with a message that says what's expected), delegate,
respond with `writeJSON`/`writeError`. Register routes in
`internal/api/server.go::New` and add the field + parameter to `Server`
and `New`, then wire it in `cmd/server/main.go`.

## 5. Tests

`internal/api/things_test.go` with a fake store (copy the `fakeNotes`
pattern): happy path, validation failures (table-driven), and the
503-when-nil case.

## 6. Frontend API module (ui/src/lib/api/things.ts)

TypeScript interface matching the Go JSON exactly, functions delegating
to `api()` from `client.ts`.

## 7. UI

Build the page/section from the primitives in
`ui/src/lib/components/ui/`. Handle errors inline, show an empty state
that says what to do, use design tokens only.

## 8. Verify

```bash
make check && make test
```

Both must pass before reporting done. If the user has the dev stack
running (`make dev` + `make services-up`), exercise the endpoint for real.
