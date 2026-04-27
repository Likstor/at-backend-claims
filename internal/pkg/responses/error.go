package responses

import (
	"context"
	"net/http"
	"strings"
)

func Error(ctx context.Context, w http.ResponseWriter, statusCode int, msg string) {
	resp := make(map[string]any)
	resp["error"] = strings.ReplaceAll(msg, "\n", "; ") // Заменяем разделитель обернутых ошибок на ;

	JSON(ctx, w, statusCode, resp)
}