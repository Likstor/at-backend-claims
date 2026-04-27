package logs

import (
	"at-backend-claims/internal/pkg/reqctx"
	"context"
	"log/slog"
)

// TODO: Добавить Setup

type handlerMiddleware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *handlerMiddleware {
	return &handlerMiddleware{next: next}
}

func (h handlerMiddleware) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h handlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(reqctx.ReqCtxKey).(reqctx.ReqCtx); ok {
		if c.CorrelationID != nil {
			rec.Add("correlation_id", *c.CorrelationID)
		}

		if c.ResponseCode != nil {
			rec.Add("response_code", *c.ResponseCode)
		}

		if c.ClaimID != nil {
			rec.Add("claim_id", *c.ClaimID)
		}
	}

	return h.next.Handle(ctx, rec)
}

func (h handlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.next.WithAttrs(attrs)
}

func (h handlerMiddleware) WithGroup(name string) slog.Handler {
	return h.next.WithGroup(name)
}
