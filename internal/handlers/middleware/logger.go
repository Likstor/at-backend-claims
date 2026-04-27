package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

var opLogger = slog.String("op", "middleware.Logger")

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(
			r.Context(),
			"request to server",
			slog.Group(
				"request",
				"method", r.Method,
				"path", r.URL.Path,
			),
		)

		start := time.Now()

		next.ServeHTTP(w, r)

		slog.InfoContext(
			r.Context(),
			"request completed",
			slog.String("request processing time", time.Since(start).String()),
		)
	})
}
