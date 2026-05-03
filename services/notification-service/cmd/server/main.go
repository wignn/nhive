package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/novelhive/notification-service/internal/events"
	grpcserver "github.com/novelhive/notification-service/internal/grpc"
	"github.com/novelhive/notification-service/internal/push"
	"github.com/novelhive/notification-service/internal/repository"
	"github.com/novelhive/pkg/logger"
	libraryv1 "github.com/novelhive/proto/library/v1"
	notificationv1 "github.com/novelhive/proto/notification/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log := logger.New("notification-service")
	defer log.Sync()

	dbURL := getEnv("DATABASE_URL", "postgres://novelhive:change_me@localhost:5432/novelhive_notifications?sslmode=disable")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal("failed to connect to NATS", zap.Error(err))
	}
	defer nc.Close()

	libraryAddr := getEnv("LIBRARY_SERVICE_ADDR", "localhost:50056")
	libConn, err := grpc.Dial(libraryAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to library service", zap.Error(err))
	}
	defer libConn.Close()
	libraryClient := libraryv1.NewLibraryServiceClient(libConn)

	repo := repository.NewPostgresNotificationRepo(pool)

	firebaseCreds := getEnv("FIREBASE_CREDENTIALS", "")
	pusher, err := push.NewFirebasePusher(firebaseCreds)
	if err != nil {
		log.Fatal("failed to initialize firebase pusher", zap.Error(err))
	}

	subscriber, err := events.NewNATSSubscriber(nc, repo, libraryClient, pusher, log)
	if err != nil {
		log.Fatal("failed to create NATS subscriber", zap.Error(err))
	}

	if err := subscriber.Start(); err != nil {
		log.Fatal("failed to start NATS subscriber", zap.Error(err))
	}

	grpcPort := getEnv("GRPC_PORT", "50057")
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()
	srv := grpcserver.NewNotificationServiceServer(repo, log)
	notificationv1.RegisterNotificationServiceServer(s, srv)

	go func() {
		log.Info("gRPC server started", zap.String("port", grpcPort))
		if err := s.Serve(lis); err != nil {
			log.Fatal("failed to serve gRPC", zap.Error(err))
		}
	}()

	log.Info("notification-service started",
		zap.String("db", dbURL),
		zap.String("nats", natsURL),
		zap.String("library_service", libraryAddr),
		zap.String("firebase_enabled", func() string {
			if firebaseCreds != "" {
				return "yes"
			}
			return "no"
		}()),
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("shutting down notification-service")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
