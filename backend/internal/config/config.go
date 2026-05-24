package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	ServerPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	SessionSecret     string
	SessionCookieName string
	SessionMaxAge     int

	CORSAllowedOrigin string

	UploadDir      string
	UploadMaxBytes int64
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env tidak ditemukan, pakai env vars dari OS")
	}

	return &Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		ServerPort: getEnv("SERVER_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "revocreative_cms"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		SessionSecret:     mustEnv("SESSION_SECRET"),
		SessionCookieName: getEnv("SESSION_COOKIE_NAME", "revocms_session"),
		SessionMaxAge:     getEnvInt("SESSION_MAX_AGE_SECONDS", 86400),

		CORSAllowedOrigin: getEnv("CORS_ALLOWED_ORIGIN", "http://localhost:5173"),

		UploadDir:      getEnv("UPLOAD_DIR", "./uploads"),
		UploadMaxBytes: int64(getEnvInt("UPLOAD_MAX_BYTES", 5242880)),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("env %s wajib di-set di .env", key)
	}
	return v
}
