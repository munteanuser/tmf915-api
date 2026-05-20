package middleware

import (
	"net/http"

	"github.com/tmf915-api/internal/handlers"
)

func APIKey(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Public endpoints
			if r.URL.Path == "/health" || r.URL.Path == "/docs" || r.URL.Path == "/swagger.json" ||
				len(r.URL.Path) > 5 && r.URL.Path[:5] == "/docs" {
				next.ServeHTTP(w, r)
				return
			}

			got := r.Header.Get("X-API-Key")
			if got == "" {
				got = r.URL.Query().Get("api_key")
			}
			if got != key {
				handlers.WriteError(w, http.StatusUnauthorized, "401", "Unauthorized", "Missing or invalid X-API-Key header")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
