package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	grpcserver "github.com/novelhive/library-service/internal/grpc"
	"github.com/novelhive/library-service/internal/repository"
	libraryv1 "github.com/novelhive/proto/library/v1"
	"github.com/novelhive/pkg/grpclog"
	"github.com/novelhive/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := logger.New("library-service")
	defer log.Sync()

	dbURL := getEnv("DATABASE_URL", "postgres://novelhive:secret@localhost:5432/novelhive_library?sslmode=disable")
	grpcPort := getEnv("GRPC_PORT", "50056")

	log.Info("connecting to database")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	log.Info("running database migrations")
	runMigrations(pool)

	libraryRepo := repository.NewPostgresLibraryRepo(pool)
	bookmarkRepo := repository.NewPostgresBookmarkRepo(pool)
	progressRepo := repository.NewPostgresProgressRepo(pool)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatal("failed to listen", zap.String("port", grpcPort), zap.Error(err))
	}

	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(grpclog.UnaryServerInterceptor(log)),
	)
	libraryv1.RegisterLibraryServiceServer(grpcSrv, grpcserver.NewLibraryServiceServer(libraryRepo, bookmarkRepo, progressRepo, log))
	reflection.Register(grpcSrv)

	log.Info("library-service started", zap.String("port", grpcPort))

	// Graceful shutdown
	go func() {
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatal("grpc serve failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down library-service")
	grpcSrv.GracefulStop()
	log.Info("library-service stopped")
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
	for _, q := range queries {
		pool.Exec(context.Background(), q)
	}
}

func getEnv(key, fb string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fb
}
