package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/natrontech/klickops-starter/internal/api"
)

// Notes implements api.NoteStore backed by PostgreSQL.
type Notes struct {
	pool *pgxpool.Pool
}

func NewNotes(pool *pgxpool.Pool) *Notes {
	return &Notes{pool: pool}
}

func (n *Notes) ListNotes(ctx context.Context) ([]api.Note, error) {
	rows, err := n.pool.Query(ctx,
		`SELECT id, text, created_at FROM notes ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes: %w", err)
	}
	notes, err := pgx.CollectRows(rows, pgx.RowToStructByPos[api.Note])
	if err != nil {
		return nil, fmt.Errorf("failed to scan notes: %w", err)
	}
	return notes, nil
}

func (n *Notes) CreateNote(ctx context.Context, text string) (api.Note, error) {
	var note api.Note
	err := n.pool.QueryRow(ctx,
		`INSERT INTO notes (text) VALUES ($1) RETURNING id, text, created_at`, text,
	).Scan(&note.ID, &note.Text, &note.CreatedAt)
	if err != nil {
		return api.Note{}, fmt.Errorf("failed to insert note: %w", err)
	}
	return note, nil
}

func (n *Notes) DeleteNote(ctx context.Context, id int64) error {
	if _, err := n.pool.Exec(ctx, `DELETE FROM notes WHERE id = $1`, id); err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}
	return nil
}
