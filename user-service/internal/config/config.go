package config

import "os"

type Config struct {
	DatabaseURL    string
	RedisURL       string
	JWTSecret      string
	GRPCPort       string
	InternalAPIKey string
}

func Load() *Config {
	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://novelhive:secret@localhost:5432/novelhive_users?sslmode=disable"),
		RedisURL:       getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:      getEnv("JWT_SECRET", "novelhive-dev-secret"),
		GRPCPort:       getEnv("GRPC_PORT", "50051"),
		InternalAPIKey: getEnv("INTERNAL_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

