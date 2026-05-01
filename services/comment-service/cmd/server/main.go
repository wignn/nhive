package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/novelhive/pkg/grpcauth"
	"github.com/novelhive/pkg/grpclog"
	"github.com/novelhive/pkg/logger"
	grpcserver "github.com/novelhive/comment-service/internal/grpc"
	"github.com/novelhive/comment-service/internal/repository"
	commentv1 "github.com/novelhive/proto/comment/v1"
	userv1 "github.com/novelhive/proto/user/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := logger.New("comment-service")
	defer log.Sync()

	dbURL := getEnv("DATABASE_URL", "postgres://novelhive:secret@localhost:5432/novelhive_comments?sslmode=disable")
	grpcPort := getEnv("GRPC_PORT", "50055")
	userServiceAddr := getEnv("USER_SERVICE_ADDR", "localhost:50051")
	apiKey := getEnv("INTERNAL_API_KEY", "")

	log.Info("connecting to database")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	log.Info("running database migrations")
	runMigrations(pool)

	// Connect to user-service for profile resolution (inject internal key)
	var userClient userv1.UserServiceClient
	userDialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(grpcauth.NewCredentials(apiKey)),
	}
	if conn, err := grpc.NewClient(userServiceAddr, userDialOpts...); err == nil {
		userClient = userv1.NewUserServiceClient(conn)
		log.Info("connected to user-service", zap.String("addr", userServiceAddr))
	} else {
		log.Warn("failed to connect to user-service (non-fatal)", zap.String("addr", userServiceAddr), zap.Error(err))
	}

	commentRepo := repository.NewPostgresCommentRepo(pool)
	likeRepo := repository.NewPostgresLikeRepo(pool)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatal("failed to listen", zap.String("port", grpcPort), zap.Error(err))
	}

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpclog.UnaryServerInterceptor(log),
			grpcauth.UnaryServerInterceptor(apiKey),
		),
	)
	commentv1.RegisterCommentServiceServer(grpcSrv, grpcserver.NewCommentServiceServer(commentRepo, likeRepo, userClient, log))
	reflection.Register(grpcSrv)

	log.Info("comment-service started", zap.String("port", grpcPort))

	// Graceful shutdown
	go func() {
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatal("grpc serve failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down comment-service")
	grpcSrv.GracefulStop()
	log.Info("comment-service stopped")
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
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fb
}
