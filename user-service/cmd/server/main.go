package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/novelhive/user-service/internal/config"
	"github.com/novelhive/user-service/internal/repository"
	"github.com/novelhive/user-service/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	// Connect to PostgreSQL
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Run migrations
	if err := runMigrations(pool); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Setup layers
	userRepo := repository.NewPostgresUserRepo(pool)
	userUC := usecase.NewUserUsecase(userRepo, cfg.JWTSecret)
	_ = userUC // Will be used when proto-gen server is ready

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	// TODO: Register proto-generated service here
	// pb.RegisterUserServiceServer(grpcServer, grpcserver.NewUserServiceServer(userUC))
	reflection.Register(grpcServer)

	log.Printf("User Service listening on :%s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func runMigrations(pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(64) PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		avatar_url TEXT DEFAULT '',
		role VARCHAR(20) DEFAULT 'reader',
		created_at TIMESTAMPTZ DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	`
	_, err := pool.Exec(context.Background(), query)
	return err
}
