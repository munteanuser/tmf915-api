package middleware

import (
	"net/http"

	"github.com/tmf915-api/internal/handlers"
)

func APIKey(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Public endpoints
			p := r.URL.Path
			if p == "/health" || p == "/docs" || p == "/swagger.json" ||
				(len(p) > 5 && p[:5] == "/docs") ||
				p == "/app" || (len(p) > 4 && p[:5] == "/app/") {
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
