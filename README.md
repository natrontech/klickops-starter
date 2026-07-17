# klickops starter

[![ci](https://github.com/natrontech/klickops-starter/actions/workflows/ci.yml/badge.svg)](https://github.com/natrontech/klickops-starter/actions/workflows/ci.yml)

A rock-solid starting point for building a web app: **SvelteKit frontend + Go
backend in a single container**, with optional **PostgreSQL**, **Valkey** and **S3
storage** wired the way [klickops](https://klickops.io) provides them.

It is a seed, not a framework. Everything in here is meant to be understood
in an afternoon and replaced by your own app. It ships with the guardrails
that make AI-assisted development work well: a `CLAUDE.md` brief, coding
rules, a feature recipe, tests, CI, and Dependabot.

*Deutsche Version: [README.de.md](README.de.md)*

## What you get

- **One container** that serves the API (`/api/*`) and the built UI. One
  Dockerfile, one deployable artifact, no orchestration puzzle.
- **Go backend** on the standard library (no framework), with a tiny
  embedded SQL migration runner and graceful degradation: the app starts
  and runs even before a database or bucket is bound.
- **SvelteKit 5 + Tailwind 4 frontend** as a static SPA, with a small
  design-token system (automatic dark mode) and a few UI primitives.
- **Working examples**: a notes CRUD (PostgreSQL), a visit counter and
  cache-aside list caching (Valkey), and file upload/download (S3). Delete
  them once your real app takes shape.
- **AI guidance built in**: `CLAUDE.md`, `.claude/rules/`, and a
  step-by-step `new-feature` skill so Claude Code, Cursor, and friends
  extend the app idiomatically instead of inventing their own patterns.
- **Quality gates**: Go and frontend tests, `make check`, GitHub Actions
  CI (format, vet, type-check, test, build, Docker build), and Dependabot
  keeping Go, npm, Actions, and base images up to date.

## Quick start

Use this repo as a GitHub template ("Use this template" button), then either
open your new repo in [GitHub Codespaces](https://github.com/features/codespaces)
(the included devcontainer preinstalls Go, Node, pnpm, Docker, and the
[Claude Code](https://claude.com/claude-code) CLI - type `claude` in the
terminal and start building, no local setup at all) or work locally:

```bash
git clone git@github.com:you/your-app.git
cd your-app
make install        # frontend dependencies (needs pnpm)
make dev            # backend on :8080, frontend on :5173
```

Open http://localhost:5173. The app runs with zero configuration; the
database and storage sections show how to enable them.

To develop with real services locally (needs Docker):

```bash
make services-up          # PostgreSQL + Valkey + S3-compatible server
cp .env.example .env      # points at those services
make dev
```

## Build it with AI

This repo is written to be extended by AI coding tools. Open the folder
with [Claude Code](https://claude.com/claude-code) (or Cursor, etc.) and
describe what you want:

> "Replace the notes demo with a customer list: name, email, a note field,
> and a CSV export."

The AI reads `CLAUDE.md` for the architecture, `.claude/rules/` for the
coding conventions, and `.claude/skills/new-feature/` for the exact recipe
to add a feature end to end (migration, handler, tests, UI). Ask it to run
`make check && make test` before it declares victory; CI enforces the same.

## Deploy on klickops

1. Push your repo to GitHub.
2. In klickops: create a project, **Deploy from repo**, pick the repo.
   klickops detects the Dockerfile, builds it, and deploys the container
   with a URL and TLS.
3. **Database**: add a PostgreSQL service to the project. klickops sees the
   `DATABASE_URL` variable declared in the Dockerfile and suggests the
   binding, accept it and the app connects and migrates on the next start.
4. **Cache**: add a Valkey (Redis-compatible) service and connect it to the
   app - the binding injects `REDIS_URL` and the visit counter plus notes
   caching light up on the next start.
5. **Storage**: add a Bucket service and bind it to the app. The app reads
   the standard AWS SDK variables the binding already injects
   (`AWS_ENDPOINT_URL_S3`, `AWS_REGION`, `AWS_ACCESS_KEY_ID`,
   `AWS_SECRET_ACCESS_KEY`, plus `S3_BUCKET`), so accept the suggested
   binding and the files API works on the next start.

No YAML, no kubectl. Any other container platform works too, the app only
needs the env vars below.

## Configuration

Everything is an environment variable with a sensible default
(`internal/config/config.go` is the single source of truth):

| Variable | Default | Purpose |
|---|---|---|
| `PORT` | `8080` | HTTP listen port |
| `DATABASE_URL` | *(unset)* | PostgreSQL connection string; unset = notes API disabled |
| `REDIS_URL` | *(unset)* | Valkey/Redis connection URL; unset = cache + visits disabled |
| `AWS_ENDPOINT_URL_S3` | *(unset)* | S3-compatible endpoint (`host:port` or URL) |
| `AWS_REGION` | `us-east-1` | S3 signing region |
| `S3_BUCKET` | *(unset)* | bucket name; unset = files API disabled |
| `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` | *(unset)* | bucket credentials |
| `UI_DIR` | `ui/build` | where the built SPA lives |

## Project structure

```
cmd/server/             Go entry point (config, wiring, graceful shutdown)
internal/config/        env var parsing, the only place os.Getenv appears
internal/api/           HTTP handlers + store interfaces + tests
internal/db/            pgx pool, embedded SQL migrations, store impls
internal/storage/       S3 implementation
ui/src/routes/          SvelteKit pages (SPA)
ui/src/lib/api/         typed fetch wrappers
ui/src/lib/components/  design-system primitives
.claude/                rules + skills that guide AI coding tools
.devcontainer/          GitHub Codespaces / devcontainer setup
```

## Commands

```bash
make help            # list all commands
make dev             # run backend + frontend for development
make check           # type-check everything
make test            # run all tests
make lint            # formatting + vet checks
make build           # production build (bin/server + ui/build)
make docker-build    # build the container image
```

## License

[MIT](LICENSE). Build something.
