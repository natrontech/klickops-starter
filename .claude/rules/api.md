# API Conventions

## Shape

- All endpoints live under `/api/`. Everything else is the SPA.
- JSON in, JSON out. Responses are objects or arrays, never bare strings.
- List endpoints return `[]`, never `null` (initialize empty slices).
- Field names are camelCase in JSON (`json:"createdAt"`).

## Errors

One error shape everywhere:

```json
{ "error": "human-readable message that says what to do about it" }
```

| Status | When |
|--------|------|
| 400 | invalid input - say which field and what's expected |
| 404 | resource not found |
| 503 | optional dependency not bound (DB/storage) - include the binding hint |
| 500 | unexpected - log details with slog, return a generic safe message |

A good 400 message tells the user how to fix the request
("text is required", "id must be a number") - never just "bad request".

## Validation

- Validate at the top of the handler, before touching any store.
- Enforce size limits: request bodies, upload sizes (`http.MaxBytesReader`),
  string lengths. Reject path-traversal characters in user-supplied keys.

## Adding an endpoint

1. Store interface method in `internal/api/` + implementation in
   `internal/db/` or `internal/storage/`.
2. Handler in `internal/api/<feature>.go`.
3. Route in `server.go::New` (`mux.HandleFunc("POST /api/things", ...)`).
4. Handler test with a fake store.
5. Frontend module in `ui/src/lib/api/<feature>.ts`.

The full recipe with code templates: `.claude/skills/new-feature/SKILL.md`.
