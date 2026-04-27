package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/novelhive/gateway/internal/config"
	"github.com/novelhive/gateway/internal/handler"
	"github.com/novelhive/gateway/internal/middleware"
	"github.com/novelhive/gateway/internal/storage"
	"github.com/novelhive/gateway/internal/store"
)

func main() {
	_ = godotenv.Load("../frontend/.env")
	_ = godotenv.Load()

	cfg := config.Load()

	r2Client, err := storage.NewR2Client(cfg)
	if err != nil {
		log.Printf("⚠️ R2 Storage not configured: %v", err)
	} else {
		log.Printf("☁️ R2 Storage initialized for bucket: %s", cfg.R2BucketName)
	}

	// Initialize in-memory data store with sample data
	dataStore := store.NewStore()
	h := handler.New(dataStore, cfg.JWTSecret, r2Client)

	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(middleware.SecurityHeaders)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rate limiter: 120 requests per minute per IP
	rl := middleware.NewRateLimiter(120)
	r.Use(rl.Middleware)

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public auth routes
		r.Post("/auth/register", h.Register)
		r.Post("/auth/login", h.Login)

		// Public data routes (optional auth for personalization)
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

		// Protected routes (require auth)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			r.Get("/auth/me", h.GetProfile)

			// Comments
			r.Post("/chapters/{chapterId}/comments", h.CreateComment)
			r.Post("/comments/{commentId}/like", h.LikeComment)

			// Library
			r.Get("/library", h.GetLibrary)
			r.Post("/library/{novelId}", h.AddToLibrary)
			r.Put("/library/{novelId}/status", h.UpdateLibraryStatus)

			// Bookmarks
			r.Get("/bookmarks", h.GetBookmarks)
			r.Post("/bookmarks", h.AddBookmark)

			// Progress
			r.Get("/progress/{novelId}", h.GetProgress)
			r.Put("/progress/{novelId}", h.SaveProgress)
		})

		// Admin routes (require auth + admin role)
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

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"novelhive-gateway"}`))
	})

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	log.Printf("🚀 NovelHive API Gateway listening on http://localhost%s", addr)
	log.Printf("📚 Admin: admin@novelhive.com / Admin123!")
	log.Printf("📖 Reader: reader@novelhive.com / Reader123!")
	log.Fatal(http.ListenAndServe(addr, r))
}
