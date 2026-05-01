package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/novelhive/pkg/grpcauth"
	"github.com/novelhive/pkg/grpclog"
	"github.com/novelhive/pkg/logger"
	grpcserver "github.com/novelhive/novel-service/internal/grpc"
	"github.com/novelhive/novel-service/internal/events"
	"github.com/novelhive/novel-service/internal/repository"
	"github.com/novelhive/novel-service/internal/usecase"
	novelv1 "github.com/novelhive/proto/novel/v1"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := logger.New("novel-service")
	defer log.Sync()

	dbURL := getEnv("DATABASE_URL", "postgres://novelhive:secret@localhost:5432/novelhive_novels?sslmode=disable")
	redisURL := getEnv("REDIS_URL", "redis://localhost:6379/1")
	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	grpcPort := getEnv("GRPC_PORT", "50052")

	log.Info("connecting to database")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	log.Info("running database migrations")
	runMigrations(pool)

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Warn("failed to parse redis URL", zap.Error(err))
	}
	rdb := redis.NewClient(opt)
	cache := repository.NewRedisCache(rdb)
	log.Info("connected to redis", zap.String("url", redisURL))

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Warn("NATS connect failed (non-fatal)", zap.String("url", natsURL), zap.Error(err))
	} else {
		log.Info("connected to NATS", zap.String("url", natsURL))
	}
	var publisher *events.NATSPublisher
	if nc != nil {
		publisher, _ = events.NewNATSPublisher(nc)
	}

	novelRepo := repository.NewPostgresNovelRepo(pool)
	chapterRepo := repository.NewPostgresChapterRepo(pool)
	genreRepo := repository.NewPostgresGenreRepo(pool)
	novelUC := usecase.NewNovelUsecase(novelRepo, chapterRepo, genreRepo, cache, publisher)

	apiKey := getEnv("INTERNAL_API_KEY", "")

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
	novelv1.RegisterNovelServiceServer(grpcSrv, grpcserver.NewNovelServiceServer(novelUC))
	reflection.Register(grpcSrv)

	log.Info("novel-service started", zap.String("port", grpcPort))

	// Graceful shutdown
	go func() {
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatal("grpc serve failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down novel-service")
	grpcSrv.GracefulStop()
	if nc != nil {
		nc.Close()
	}
	log.Info("novel-service stopped")
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
