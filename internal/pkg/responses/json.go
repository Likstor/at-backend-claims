package responses

import (
	"at-backend-claims/internal/pkg/reqctx"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

// JSON формирует ответ в формате JSON и записывает его в http.ResponseWriter.
func JSON(ctx context.Context, w http.ResponseWriter, statusCode int, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())

		InternalServerError(ctx, w)
		return
	}

	correlationID, ok := reqctx.GetCorrelationID(ctx)
	if ok {
		w.Header().Set("X-Correlation-ID", correlationID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(bytes); err != nil {
		slog.ErrorContext(ctx, err.Error())

		return
	}

	ctx = reqctx.WithResponseCode(ctx, statusCode)

	slog.InfoContext(ctx, "written answer")
}
