package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	dbURL := getEnv("DATABASE_URL", "postgres://novelhive:secret@localhost:5432/novelhive_comments?sslmode=disable")
	grpcPort := getEnv("GRPC_PORT", "50055")

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("DB connect failed: %v", err)
	}
	defer pool.Close()
	runMigrations(pool)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Listen failed: %v", err)
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	log.Printf("Comment Service listening on :%s", grpcPort)
	grpcServer.Serve(lis)
}

func runMigrations(pool *pgxpool.Pool) {
	queries := []string{
		`CREATE EXTENSION IF NOT EXISTS ltree`,
		`CREATE TABLE IF NOT EXISTS comments (
			id VARCHAR(64) PRIMARY KEY, chapter_id VARCHAR(64) NOT NULL,
			user_id VARCHAR(64) NOT NULL, username VARCHAR(50) DEFAULT '',
			avatar_url TEXT DEFAULT '', content TEXT NOT NULL,
			parent_id VARCHAR(64) REFERENCES comments(id),
			path ltree NOT NULL, likes_count INT DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW(), deleted_at TIMESTAMPTZ
		)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_chapter ON comments(chapter_id)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_path ON comments USING GIST (path)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_parent ON comments(parent_id)`,
		`CREATE TABLE IF NOT EXISTS comment_likes (
			comment_id VARCHAR(64) REFERENCES comments(id) ON DELETE CASCADE,
			user_id VARCHAR(64) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			PRIMARY KEY(comment_id, user_id)
		)`,
	}
	for _, q := range queries {
		pool.Exec(context.Background(), q)
	}
}

func getEnv(key, fb string) string {
	if v := os.Getenv(key); v != "" { return v }
	return fb
}
