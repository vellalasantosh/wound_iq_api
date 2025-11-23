package config

import (
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    DB_DSN   string
    Port     string
    LogLevel string
}

func Load() *Config {
    _ = godotenv.Load()
    return &Config{
        DB_DSN:   getenv("DB_DSN", "postgres://postgres:postgres@localhost:5432/wound_iq?sslmode=disable"),
        Port:     getenv("PORT", "8080"),
        LogLevel: getenv("LOG_LEVEL", "info"),
    }
}

func getenv(k, fallback string) string {
    v := os.Getenv(k)
    if v == "" {
        return fallback
    }
    return v
}
