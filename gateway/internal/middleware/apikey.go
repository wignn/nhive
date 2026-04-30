package middleware

import (
	"net/http"
)

// APIKeyMiddleware validates the X-Internal-Key header on every request.
// This key is injected server-side by the Next.js proxy — the browser never sees it.
// If apiKey is empty (dev mode), the check is skipped and all requests pass through.
func APIKeyMiddleware(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if apiKey == "" {
				// No key configured — skip validation (dev mode)
				next.ServeHTTP(w, r)
				return
			}

			key := r.Header.Get("X-Internal-Key")
			if key != apiKey {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"unauthorized: invalid or missing API key"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
