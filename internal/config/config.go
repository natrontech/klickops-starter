// Package config loads all runtime configuration from environment variables.
// Every value has a sensible default; the app must start with zero config.
package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	// RedisURL is what a klickops Valkey binding injects (REDIS_URL) —
	// any Redis-compatible server works.
	RedisURL string
	UIDir    string
	S3       S3
}

// S3 holds object-storage credentials. The env names are the AWS-SDK
// convention, which is exactly what a klickops Bucket binding injects —
// binding a Bucket service to the app needs no further configuration.
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
		RedisURL:    os.Getenv("REDIS_URL"),
		UIDir:       getenv("UI_DIR", "ui/build"),
		S3: S3{
			Endpoint:  os.Getenv("AWS_ENDPOINT_URL_S3"),
			Region:    getenv("AWS_REGION", "us-east-1"),
			Bucket:    os.Getenv("S3_BUCKET"),
			AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
