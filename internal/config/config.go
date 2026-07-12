// Package config loads all runtime configuration from environment variables.
// Every value has a sensible default; the app must start with zero config.
package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	UIDir       string
	S3          S3
}

// S3 holds object-storage credentials. On klickops these are injected by
// binding a Bucket service to the app (see README "Deploy on klickops").
type S3 struct {
	Endpoint  string
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
}

func (s S3) Enabled() bool { return s.Bucket != "" }

func Load() Config {
	return Config{
		Port:        getenv("PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		UIDir:       getenv("UI_DIR", "ui/build"),
		S3: S3{
			Endpoint:  os.Getenv("S3_ENDPOINT"),
			Region:    getenv("S3_REGION", "us-east-1"),
			Bucket:    os.Getenv("S3_BUCKET"),
			AccessKey: os.Getenv("S3_ACCESS_KEY"),
			SecretKey: os.Getenv("S3_SECRET_KEY"),
		},
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
