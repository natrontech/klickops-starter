# Go Conventions

## Style

- Idiomatic Go. Return early, no deep nesting. Boring beats clever.
- Imports in three groups: stdlib, external, internal.
- Exported PascalCase, unexported camelCase. No `Get` prefixes.
- `gofmt` clean at all times (`make format`).

## Errors

```go
if err != nil {
    return fmt.Errorf("failed to <action>: %w", err)
}
```

Always wrap with context via `%w`. Log internal details with slog, return a
safe human message to the API caller - never `err.Error()` in a response.

## Logging

`log/slog` with key-value pairs: `slog.Info("applied migration", "version", v)`.
Never `fmt.Println` or `log.Printf`.

## Configuration

All env vars are parsed in `internal/config/config.go` and nowhere else.
Every value has a default; the app must start with zero configuration.
New env var → add it to `config.go`, `.env.example`, and the Dockerfile
`ENV` block in the same change.

## HTTP handlers

- Stdlib `net/http` with method+path patterns (`"GET /api/notes/{id}"`).
- Handlers are thin: validate input → call the store interface → respond
  with `writeJSON`/`writeError`.
- Stores are interfaces defined in `internal/api` (consumption point),
  implemented in `internal/db` / `internal/storage`.
- Optional dependencies (DB, storage) may be nil - guard with a 503 + a
  hint that tells the user how to bind the service.

## Testing

- Every handler and helper gets a test in the same package (`_test.go`).
- Table-driven tests for multiple scenarios.
- Mock at the boundary only: fake store structs implementing the interface
  (see `internal/api/notes_test.go`). Never mock in-process logic, never
  require a live database in unit tests.
- `httptest` for handler tests, exercising the real router via `api.New`.

## Database

- Plain SQL through pgx. No ORM.
- Schema changes are new numbered files in `internal/db/migrations/` -
  applied in filename order at startup, never edit an applied migration.
- Always parameterized queries (`$1`), never string-built SQL.
