package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/novelhive/novel-service/internal/events"
	"github.com/novelhive/novel-service/internal/repository"
	"github.com/novelhive/novel-service/internal/usecase"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"os"
)

func main() {
	dbURL := getEnv("DATABASE_URL", "postgres://novelhive:secret@localhost:5432/novelhive_novels?sslmode=disable")
	redisURL := getEnv("REDIS_URL", "redis://localhost:6379/1")
	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	grpcPort := getEnv("GRPC_PORT", "50052")

	// PostgreSQL
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("DB connect failed: %v", err)
	}
	defer pool.Close()
	runMigrations(pool)

	// Redis
	opt, _ := redis.ParseURL(redisURL)
	rdb := redis.NewClient(opt)
	cache := repository.NewRedisCache(rdb)

	// NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Printf("NATS connect failed (non-fatal): %v", err)
	}
	var publisher *events.NATSPublisher
	if nc != nil {
		publisher, _ = events.NewNATSPublisher(nc)
	}

	// Repos & usecase
	novelRepo := repository.NewPostgresNovelRepo(pool)
	chapterRepo := repository.NewPostgresChapterRepo(pool)
	genreRepo := repository.NewPostgresGenreRepo(pool)
	_ = usecase.NewNovelUsecase(novelRepo, chapterRepo, genreRepo, cache, publisher)

	// gRPC
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Listen failed: %v", err)
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	log.Printf("Novel Service listening on :%s", grpcPort)
	grpcServer.Serve(lis)
}

func runMigrations(pool *pgxpool.Pool) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS novels (
			id VARCHAR(64) PRIMARY KEY, title VARCHAR(500) NOT NULL, slug VARCHAR(500) UNIQUE NOT NULL,
			synopsis TEXT DEFAULT '', cover_url TEXT DEFAULT '', author VARCHAR(255) DEFAULT '',
			status VARCHAR(50) DEFAULT 'ongoing', total_chapters INT DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW(), updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS genres (
			id SERIAL PRIMARY KEY, name VARCHAR(100) UNIQUE NOT NULL, slug VARCHAR(100) UNIQUE NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS novel_genres (
			novel_id VARCHAR(64) REFERENCES novels(id) ON DELETE CASCADE,
			genre_id INT REFERENCES genres(id) ON DELETE CASCADE,
			PRIMARY KEY (novel_id, genre_id)
		)`,
		`CREATE TABLE IF NOT EXISTS chapters (
			id VARCHAR(64) PRIMARY KEY, novel_id VARCHAR(64) NOT NULL REFERENCES novels(id) ON DELETE CASCADE,
			number INT NOT NULL, title VARCHAR(500) DEFAULT '', content TEXT NOT NULL,
			word_count INT DEFAULT 0, created_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(novel_id, number)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_novels_slug ON novels(slug)`,
		`CREATE INDEX IF NOT EXISTS idx_chapters_novel ON chapters(novel_id, number)`,
		// Seed genres
		`INSERT INTO genres (name, slug) VALUES
			('Fantasy', 'fantasy'), ('Action', 'action'), ('Romance', 'romance'),
			('Adventure', 'adventure'), ('Sci-Fi', 'sci-fi'), ('Mystery', 'mystery'),
			('Horror', 'horror'), ('Comedy', 'comedy'), ('Drama', 'drama'),
			('Slice of Life', 'slice-of-life'), ('Martial Arts', 'martial-arts'),
			('Isekai', 'isekai'), ('Wuxia', 'wuxia'), ('Xianxia', 'xianxia')
		ON CONFLICT DO NOTHING`,
	}
	for _, q := range queries {
		pool.Exec(context.Background(), q)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
