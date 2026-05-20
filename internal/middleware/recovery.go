package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog"
	"github.com/tmf915-api/internal/handlers"
)

func Recovery(log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error().
						Interface("panic", err).
						Str("stack", string(debug.Stack())).
						Msg("recovered from panic")
					handlers.WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", "An unexpected error occurred")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
