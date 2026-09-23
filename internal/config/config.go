package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			getEnv("PGUSER", "program"),
			getEnv("PGPASSWORD", "test"),
			getEnv("PGHOST", "localhost"),
			getEnv("PGPORT", "5432"),
			getEnv("PGDATABASE", "persons"),
		)
	}

	return Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: dbURL,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
