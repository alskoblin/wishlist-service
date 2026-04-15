package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr       string
	DatabaseURL    string
	JWTSecret      string
	JWTTTL         time.Duration
	MigrationsPath string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	ttlRaw := getEnv("JWT_TTL", "24h")
	ttl, err := time.ParseDuration(ttlRaw)
	if err != nil {
		return nil, fmt.Errorf("parse JWT_TTL: %w", err)
	}

	cfg := &Config{
		HTTPAddr:       getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://wishlist:wishlist@postgres:5432/wishlist?sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "super-secret-key"),
		JWTTTL:         ttl,
		MigrationsPath: getEnv("MIGRATIONS_PATH", "file://migrations"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
