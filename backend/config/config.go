package config

import (
	"os"
)

type Config struct {
	DatabaseURL    string
	FiveFiltersURL string
	Port           string
}

func Load() *Config {
	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@postgres:5432"),
		FiveFiltersURL: getEnv("FIVEFILTERS_URL", "http://fivefilters-service:8081"),
		Port:           getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}