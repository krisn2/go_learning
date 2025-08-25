package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL     string
	JWTSecret       string
	AdminListSecret string
	Port            string
}

func Load() *Config {
	cfg := &Config{
		DatabaseURL:     getenv("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=go_social port=5432 sslmode=disable"),
		JWTSecret:       getenv("JWT_SECRET", "dev_super_secret_change_me"),
		AdminListSecret: getenv("ADMIN_LIST_SECRET", "let_me_list_users"),
		Port:            getenv("PORT", "8080"),
	}

	// Validate critical config
	if len(cfg.JWTSecret) < 32 {
		log.Println("WARNING: JWT_SECRET should be at least 32 characters for security")
	}

	if cfg.JWTSecret == "dev_super_secret_change_me" {
		log.Println("WARNING: Using default JWT secret - change this in production!")
	}

	return cfg
}

func getenv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
