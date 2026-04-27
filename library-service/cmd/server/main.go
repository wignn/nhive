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
	dbURL := getEnv("DATABASE_URL", "postgres://novelhive:secret@localhost:5432/novelhive_library?sslmode=disable")
	grpcPort := getEnv("GRPC_PORT", "50056")

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil { log.Fatalf("DB connect failed: %v", err) }
	defer pool.Close()
	runMigrations(pool)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil { log.Fatalf("Listen failed: %v", err) }

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	log.Printf("Library Service listening on :%s", grpcPort)
	grpcServer.Serve(lis)
}

func runMigrations(pool *pgxpool.Pool) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS reading_lists (
			id VARCHAR(64) PRIMARY KEY, user_id VARCHAR(64) NOT NULL,
			novel_id VARCHAR(64) NOT NULL, status VARCHAR(20) DEFAULT 'reading',
			created_at TIMESTAMPTZ DEFAULT NOW(), UNIQUE(user_id, novel_id)
		)`,
		`CREATE TABLE IF NOT EXISTS bookmarks (
			id VARCHAR(64) PRIMARY KEY, user_id VARCHAR(64) NOT NULL,
			novel_id VARCHAR(64) NOT NULL, chapter_id VARCHAR(64) NOT NULL,
			note TEXT DEFAULT '', created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS reading_progress (
			user_id VARCHAR(64) NOT NULL, novel_id VARCHAR(64) NOT NULL,
			chapter_number INT NOT NULL DEFAULT 1, scroll_position FLOAT DEFAULT 0,
			updated_at TIMESTAMPTZ DEFAULT NOW(), PRIMARY KEY(user_id, novel_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_reading_lists_user ON reading_lists(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_bookmarks_user ON bookmarks(user_id)`,
	}
	for _, q := range queries { pool.Exec(context.Background(), q) }
}

func getEnv(key, fb string) string {
	if v := os.Getenv(key); v != "" { return v }
	return fb
}
