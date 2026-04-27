package middleware

import (
	"at-backend-claims/internal/pkg/reqctx"
	"net/http"

	"github.com/google/uuid"
)

func Correlation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.NewString()
		}

		r = r.WithContext(reqctx.WithCorrelationID(r.Context(), correlationID))

		next.ServeHTTP(w, r)
	})
}
