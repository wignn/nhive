package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/novelhive/gateway/internal/clients"
	"github.com/novelhive/gateway/internal/config"
	"github.com/novelhive/gateway/internal/handler"
	"github.com/novelhive/gateway/internal/middleware"
	"github.com/novelhive/gateway/internal/storage"
	"github.com/novelhive/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	log := logger.New("gateway")
	defer log.Sync()

	_ = godotenv.Load("../frontend/.env")
	_ = godotenv.Load()

	cfg := config.Load()

	// R2 image storage
	r2Client, err := storage.NewR2Client(cfg)
	if err != nil {
		log.Warn("R2 Storage not configured", zap.Error(err))
	} else {
		log.Info("R2 Storage initialized", zap.String("bucket", cfg.R2BucketName))
	}

	// gRPC clients to all microservices
	log.Info("connecting to microservices",
		zap.String("user_service", cfg.UserServiceAddr),
		zap.String("novel_service", cfg.NovelServiceAddr),
		zap.String("comment_service", cfg.CommentServiceAddr),
		zap.String("library_service", cfg.LibraryServiceAddr),
	)
	svcClients := clients.New(
		cfg.UserServiceAddr,
		cfg.NovelServiceAddr,
		cfg.CommentServiceAddr,
		cfg.LibraryServiceAddr,
	)
	defer svcClients.Close()

	h := handler.New(svcClients, cfg.JWTSecret, r2Client)

	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(middleware.SecurityHeaders)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173", "https://*.novelhive.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	rl := middleware.NewRateLimiter(200)
	r.Use(rl.Middleware)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", h.Register)
		r.Post("/auth/login", h.Login)
		
		r.Group(func(r chi.Router) {
			r.Use(middleware.OptionalAuth(cfg.JWTSecret))
			r.Get("/novels", h.ListNovels)
			r.Get("/novels/{slug}", h.GetNovel)
			r.Get("/novels/{slug}/chapters", h.ListChapters)
			r.Get("/novels/{slug}/chapters/{number}", h.ReadChapter)
			r.Get("/search", h.Search)
			r.Get("/search/autocomplete", h.Autocomplete)
			r.Get("/chapters/{chapterId}/comments", h.ListComments)
			r.Get("/genres", h.ListGenres)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			r.Get("/auth/me", h.GetProfile)

			r.Post("/chapters/{chapterId}/comments", h.CreateComment)
			r.Post("/comments/{commentId}/like", h.LikeComment)

			r.Get("/library", h.GetLibrary)
			r.Post("/library/{novelId}", h.AddToLibrary)
			r.Delete("/library/{novelId}", h.RemoveFromLibrary)
			r.Put("/library/{novelId}/status", h.UpdateLibraryStatus)

			r.Get("/bookmarks", h.GetBookmarks)
			r.Post("/bookmarks", h.AddBookmark)

			r.Get("/progress/{novelId}", h.GetProgress)
			r.Put("/progress/{novelId}", h.SaveProgress)
		})

		// Admin routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			r.Use(middleware.AdminOnly)

			r.Get("/admin/novels", h.AdminListNovels)
			r.Post("/admin/novels", h.AdminCreateNovel)
			r.Post("/admin/upload", h.AdminUploadImage)
			r.Put("/admin/novels/{id}", h.AdminUpdateNovel)
			r.Delete("/admin/novels/{id}", h.AdminDeleteNovel)

			r.Get("/admin/chapters", h.AdminListChapters)
			r.Post("/admin/chapters", h.AdminCreateChapter)
			r.Put("/admin/chapters/{id}", h.AdminUpdateChapter)
			r.Delete("/admin/chapters/{id}", h.AdminDeleteChapter)

			r.Get("/admin/users", h.AdminListUsers)
			r.Put("/admin/users/{id}/role", h.AdminUpdateUserRole)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"novelhive-gateway","version":"2.0"}`))
	})

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	log.Info("gateway started", zap.String("addr", "http://localhost"+addr))

	// Graceful shutdown
	srv := &http.Server{Addr: addr, Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("http server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down gateway")
	srv.Close()
	log.Info("gateway stopped")
}
