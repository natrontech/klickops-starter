.PHONY: help install dev dev-backend dev-ui build check test lint format docker-build services-up services-down

help: ## Show available commands
	@grep -E '^[a-z-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

install: ## Install frontend dependencies
	cd ui && pnpm install

dev: ## Run backend (:8080) and frontend (:5173) together
	@trap 'kill 0' INT TERM; $(MAKE) dev-backend & $(MAKE) dev-ui & wait

dev-backend: ## Run the Go backend with .env loaded (:8080)
	@set -a; [ -f .env ] && . ./.env; set +a; go run ./cmd/server

dev-ui: ## Run the SvelteKit dev server (:5173, proxies /api to :8080)
	cd ui && pnpm dev

build: ## Build the UI and the Go binary
	cd ui && pnpm build
	go build -o bin/server ./cmd/server

check: ## Type-check Go and frontend
	go vet ./...
	cd ui && pnpm check

test: ## Run all tests
	go test ./...
	cd ui && pnpm test

lint: ## Check formatting
	@test -z "$$(gofmt -l .)" || (gofmt -l . && echo "run: make format" && exit 1)
	go vet ./...
	cd ui && pnpm lint

format: ## Format all code
	gofmt -w .
	cd ui && pnpm format

docker-build: ## Build the production container image
	docker build -t klickops-starter .

services-up: ## Start local PostgreSQL + S3 (docker compose)
	docker compose up -d

services-down: ## Stop local services
	docker compose down
