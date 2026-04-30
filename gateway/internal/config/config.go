package config

import "os"

type Config struct {
	HTTPPort      string
	JWTSecret     string
	InternalAPIKey string

	// gRPC service addresses
	UserServiceAddr    string
	NovelServiceAddr   string
	CommentServiceAddr string
	LibraryServiceAddr string

	// Cloudflare R2
	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2BucketName      string
	R2PublicURL       string
	R2Endpoint        string
}

func Load() *Config {
	return &Config{
		HTTPPort:       getEnv("HTTP_PORT", "8080"),
		JWTSecret:      getEnv("JWT_SECRET", "novelhive-dev-secret"),
		InternalAPIKey: getEnv("INTERNAL_API_KEY", ""),

		UserServiceAddr:    getEnv("USER_SERVICE_ADDR", "localhost:50051"),
		NovelServiceAddr:   getEnv("NOVEL_SERVICE_ADDR", "localhost:50052"),
		CommentServiceAddr: getEnv("COMMENT_SERVICE_ADDR", "localhost:50055"),
		LibraryServiceAddr: getEnv("LIBRARY_SERVICE_ADDR", "localhost:50056"),

		R2AccountID:       getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
		R2SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
		R2BucketName:      getEnv("R2_BUCKET_NAME", ""),
		R2PublicURL:       getEnv("R2_PUBLIC_URL", ""),
		R2Endpoint:        getEnv("R2_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

