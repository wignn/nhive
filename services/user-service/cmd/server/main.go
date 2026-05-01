package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/novelhive/user-service/internal/config"
	grpcserver "github.com/novelhive/user-service/internal/grpc"
	"github.com/novelhive/user-service/internal/repository"
	"github.com/novelhive/user-service/internal/usecase"
	"github.com/novelhive/pkg/grpcauth"
	"github.com/novelhive/pkg/grpclog"
	"github.com/novelhive/pkg/logger"
	userv1 "github.com/novelhive/proto/user/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := logger.New("user-service")
	defer log.Sync()

	cfg := config.Load()

	log.Info("connecting to database")
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	log.Info("running database migrations")
	if err := runMigrations(pool); err != nil {
		log.Fatal("failed to run migrations", zap.Error(err))
	}

	userRepo := repository.NewPostgresUserRepo(pool)
	userUC := usecase.NewUserUsecase(userRepo, cfg.JWTSecret)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatal("failed to listen", zap.String("port", cfg.GRPCPort), zap.Error(err))
	}

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpclog.UnaryServerInterceptor(log),
			grpcauth.UnaryServerInterceptor(cfg.InternalAPIKey),
		),
	)
	userv1.RegisterUserServiceServer(grpcSrv, grpcserver.NewUserServiceServer(userUC))
	reflection.Register(grpcSrv)

	log.Info("user-service started", zap.String("port", cfg.GRPCPort))

	go func() {
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatal("grpc serve failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down user-service")
	grpcSrv.GracefulStop()
	log.Info("user-service stopped")
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
