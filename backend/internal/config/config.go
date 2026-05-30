package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	UploadDir   string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://goodcity:goodcity@localhost:5432/goodcity?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret"),
		Port:        port,
		UploadDir:   getEnv("UPLOAD_DIR", "./uploads"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
