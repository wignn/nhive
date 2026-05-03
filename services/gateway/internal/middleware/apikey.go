package middleware

import (
	"crypto/subtle"
	"net/http"
)

func APIKeyMiddleware(apiKeys []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(apiKeys) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			key := r.Header.Get("X-API-Key")
			if !hasValidAPIKey(key, apiKeys) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"unauthorized: invalid or missing API key"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func hasValidAPIKey(key string, apiKeys []string) bool {
	if key == "" {
		return false
	}

	for _, apiKey := range apiKeys {
		if subtle.ConstantTimeCompare([]byte(key), []byte(apiKey)) == 1 {
			return true
		}
	}
	return false
}
