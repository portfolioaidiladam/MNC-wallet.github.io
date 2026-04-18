// Package config memuat environment variable aplikasi.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config merepresentasikan konfigurasi runtime aplikasi yang di-load
// dari environment variable (optional: .env file di CWD).
type Config struct {
	AppEnv   string
	HTTPPort string
	TZ       string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	RedisAddr     string
	RedisPassword string

	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration

	AsynqConcurrency int
}

// Load membaca konfigurasi dari env (dan optional .env). Return error
// kalau required variable kosong.
func Load() (*Config, error) {
	_ = godotenv.Load()

	accessTTLMin, _ := strconv.Atoi(getenv("JWT_ACCESS_TTL_MINUTES", "15"))
	refreshTTLHr, _ := strconv.Atoi(getenv("JWT_REFRESH_TTL_HOURS", "168"))
	concurrency, _ := strconv.Atoi(getenv("ASYNQ_CONCURRENCY", "10"))

	cfg := &Config{
		AppEnv:   getenv("APP_ENV", "development"),
		HTTPPort: getenv("HTTP_PORT", "8000"),
		TZ:       getenv("TZ", "Asia/Jakarta"),

		DBHost:     getenv("POSTGRES_HOST", "localhost"),
		DBPort:     getenv("POSTGRES_PORT", "5432"),
		DBUser:     getenv("POSTGRES_USER", "mncwallet"),
		DBPassword: getenv("POSTGRES_PASSWORD", "mncwallet"),
		DBName:     getenv("POSTGRES_DB", "mncwallet"),
		DBSSLMode:  getenv("POSTGRES_SSLMODE", "disable"),

		RedisAddr:     getenv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),

		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTAccessTTL:  time.Duration(accessTTLMin) * time.Minute,
		JWTRefreshTTL: time.Duration(refreshTTLHr) * time.Hour,

		AsynqConcurrency: concurrency,
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("config: JWT_SECRET is required")
	}
	return cfg, nil
}

// DatabaseDSN membangun DSN Postgres untuk GORM.
func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode, c.TZ,
	)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}