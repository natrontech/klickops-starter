// Package db is the PostgreSQL layer: connection, migrations, and the
// store implementations behind the interfaces in internal/api.
package db

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Connect opens a pool and waits for the database to accept connections.
// The retry loop covers the common case of the app starting before the
// database is ready.
func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DATABASE_URL: %w", err)
	}
	for attempt := 1; ; attempt++ {
		if err = pool.Ping(ctx); err == nil {
			return pool, nil
		}
		if attempt >= 10 {
			pool.Close()
			return nil, fmt.Errorf("database not reachable after %d attempts: %w", attempt, err)
		}
		slog.Info("waiting for database", "attempt", attempt)
		time.Sleep(time.Second)
	}
}

// Migrate applies every file in migrations/ (sorted by name) that has not
// been applied yet, tracked in the schema_migrations table. Add a new
// numbered .sql file to change the schema; never edit an applied one.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`); err != nil {
		return fmt.Errorf("failed to create schema_migrations: %w", err)
	}

	entries, err := migrations.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

	for _, entry := range entries {
		version := entry.Name()
		var applied bool
		if err := pool.QueryRow(ctx,
			`SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = $1)`, version,
		).Scan(&applied); err != nil {
			return fmt.Errorf("failed to check migration %s: %w", version, err)
		}
		if applied {
			continue
		}
		sql, err := migrations.ReadFile("migrations/" + version)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", version, err)
		}
		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin migration %s: %w", version, err)
		}
		if _, err := tx.Exec(ctx, string(sql)); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("migration %s failed: %w", version, err)
		}
		if _, err := tx.Exec(ctx,
			`INSERT INTO schema_migrations (version) VALUES ($1)`, version); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to record migration %s: %w", version, err)
		}
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", version, err)
		}
		slog.Info("applied migration", "version", version)
	}
	return nil
}
