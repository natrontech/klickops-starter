package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/natrontech/klickops-starter/internal/api"
	"github.com/natrontech/klickops-starter/internal/config"
	"github.com/natrontech/klickops-starter/internal/db"
	"github.com/natrontech/klickops-starter/internal/storage"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	cfg := config.Load()
	ctx := context.Background()

	var notes api.NoteStore
	if cfg.DatabaseURL != "" {
		pool, err := db.Connect(ctx, cfg.DatabaseURL)
		if err != nil {
			slog.Error("failed to connect to database", "error", err)
			os.Exit(1)
		}
		defer pool.Close()
		if err := db.Migrate(ctx, pool); err != nil {
			slog.Error("failed to run migrations", "error", err)
			os.Exit(1)
		}
		notes = db.NewNotes(pool)
		slog.Info("database connected")
	} else {
		slog.Info("DATABASE_URL not set - notes API disabled")
	}

	var blobs api.BlobStore
	if cfg.S3.Enabled() {
		store, err := storage.NewS3(ctx, cfg.S3)
		if err != nil {
			slog.Error("failed to set up storage", "error", err)
			os.Exit(1)
		}
		blobs = store
		slog.Info("storage connected", "bucket", cfg.S3.Bucket)
	} else {
		slog.Info("S3_BUCKET not set - files API disabled")
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           api.New(notes, blobs, cfg.UIDir),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown failed", "error", err)
	}
	slog.Info("shut down cleanly")
}
