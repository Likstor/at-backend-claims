package logs

import (
	"at-backend-claims/internal/pkg/apperror"
	"context"
	"log/slog"
)

func Error(ctx context.Context, err error) {
	slog.ErrorContext(apperror.GetErrorCtx(ctx, err), err.Error())
}
