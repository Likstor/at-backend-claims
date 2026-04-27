package health

import (
	"log/slog"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte(http.StatusText(http.StatusOK))); err != nil {
		slog.WarnContext(r.Context(), err.Error())
	}
}

func Setup(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", health)
}
