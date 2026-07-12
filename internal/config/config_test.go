package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	cfg := Load()
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want 8080", cfg.Port)
	}
	if cfg.UIDir != "ui/build" {
		t.Errorf("UIDir = %q, want ui/build", cfg.UIDir)
	}
	if cfg.S3.Region != "us-east-1" {
		t.Errorf("S3.Region = %q, want us-east-1", cfg.S3.Region)
	}
	if cfg.S3.Enabled() {
		t.Error("S3.Enabled() = true without S3_BUCKET")
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("PORT", "9000")
	t.Setenv("DATABASE_URL", "postgres://u:p@localhost/db")
	t.Setenv("S3_BUCKET", "uploads")

	cfg := Load()
	if cfg.Port != "9000" {
		t.Errorf("Port = %q, want 9000", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://u:p@localhost/db" {
		t.Errorf("DatabaseURL = %q", cfg.DatabaseURL)
	}
	if !cfg.S3.Enabled() {
		t.Error("S3.Enabled() = false with S3_BUCKET set")
	}
}
