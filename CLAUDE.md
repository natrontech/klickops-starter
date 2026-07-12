# klickops-starter

A single-container full-stack web app: Go backend + SvelteKit frontend, with
optional PostgreSQL and S3 storage. This is a seed for the user's own app -
the notes/files demo exists to show the patterns and is MEANT to be replaced
by whatever the user asks you to build.

## Architecture (do not change this shape)

- **One container.** The Go binary serves the API under `/api/*` and the
  built SvelteKit SPA for everything else. No second service, no separate
  frontend container, no SSR/Node runtime in production.
- **Frontend** (`ui/`): SvelteKit 5 (runes) + Tailwind 4, `adapter-static`
  with `index.html` fallback. Pure SPA - keep `ssr = false`.
- **Backend** (`cmd/server`, `internal/`): Go stdlib `net/http` with 1.22+
  route patterns. No web framework.
- **Database**: PostgreSQL via `DATABASE_URL`, optional. When unset, DB
  endpoints return 503 with a helpful hint - never panic, never require it.
- **Storage**: any S3-compatible endpoint via `S3_*` env vars, optional,
  same graceful degradation.
- Deployed on [klickops](https://klickops.io): push to GitHub, connect the
  repo, klickops builds the Dockerfile and injects `DATABASE_URL` / `S3_*`
  by binding services. Keep the Dockerfile's `EXPOSE` and `ENV` lines - the
  platform reads them to suggest bindings.

## Commands

```bash
make install         # frontend deps (pnpm)
make dev             # backend :8080 + frontend :5173 (vite proxies /api)
make services-up     # local PostgreSQL + S3 via docker compose
make check           # go vet + svelte-check   - run after every change
make test            # go test + vitest        - run after every change
make lint            # gofmt + vet + prettier check
make format          # fix formatting
make build           # UI build + Go binary (bin/server)
make docker-build    # production image
```

Before telling the user something is done: `make check && make test` must
pass. If you changed formatting-sensitive files, `make lint` too.

## Package map

```
cmd/server/            → main: config, wiring, HTTP server, shutdown
internal/config/       → all env var parsing (only place os.Getenv appears)
internal/api/          → HTTP handlers + store interfaces (NoteStore, BlobStore)
internal/db/           → pgx pool, migrations (embedded .sql), store impls
internal/db/migrations → numbered .sql files, applied in order at startup
internal/storage/      → S3 implementation of BlobStore
ui/src/routes/         → SvelteKit pages
ui/src/lib/api/        → typed fetch wrappers (always via client.ts::api)
ui/src/lib/components/ui/ → design-system primitives (Button, Card, Input, Badge)
ui/src/lib/utils/      → shared frontend helpers
```

## How to add a feature (the recipe)

Full walkthrough: `.claude/skills/new-feature/SKILL.md`. Short version:

1. Migration: new numbered file in `internal/db/migrations/` (never edit an
   applied one).
2. Store: interface method in `internal/api/`, pgx implementation in
   `internal/db/`.
3. Handler: `internal/api/<feature>.go` - validate, delegate, respond via
   `writeJSON`/`writeError`. Register the route in `server.go::New`.
4. Handler test with a fake store (see `notes_test.go`).
5. Frontend API module: `ui/src/lib/api/<feature>.ts` using `api()`.
6. Page/component using the design-system primitives.

## Conventions

- Go: `.claude/rules/go.md` - return early, wrap errors with `%w`, slog,
  table-driven tests, mock at boundaries only.
- Frontend: `.claude/rules/frontend.md` - Svelte 5 runes, design tokens
  only (no hardcoded colors), pnpm, all fetches through `lib/api/client.ts`.
- API: `.claude/rules/api.md` - error shape `{"error": "..."}`, status code
  rules, validation at the boundary.
- Design: `.claude/rules/design.md` - tokens in `app.css`, component
  patterns, dark mode is automatic.
- Git: `.claude/rules/git.md` - Conventional Commits, commit per coherent
  unit once `make check && make test` pass, never batch unrelated changes.

## What NOT to do

- Don't add a web framework (gin/echo/fiber) - stdlib routing is enough.
- Don't add an ORM - pgx with plain SQL is the pattern.
- Don't split into multiple containers or add SSR.
- Don't store secrets in code; everything configurable is an env var with a
  default in `internal/config`.
- Don't call `fetch` directly in components - extend `ui/src/lib/api/`.
- Don't hardcode colors - use the tokens in `ui/src/app.css`.
- The notes/files demo is disposable: when the user's real app takes shape,
  delete what they don't need (including this paragraph's feature files).
