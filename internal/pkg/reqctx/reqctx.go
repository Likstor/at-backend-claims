package reqctx

import (
	"at-backend-claims/internal/pkg/roles"
	"context"

	"github.com/google/uuid"
)

type keyType int

const ReqCtxKey = keyType(0)

type ReqCtx struct {
	CorrelationID *string
	UserID        *uuid.UUID
	ResponseCode  *int
	ClaimID       *uint64
	Role          *roles.Role
}

func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok {
		c.CorrelationID = &correlationID
		return context.WithValue(ctx, ReqCtxKey, c)
	}

	return context.WithValue(ctx, ReqCtxKey, ReqCtx{CorrelationID: &correlationID})
}

func GetCorrelationID(ctx context.Context) (string, bool) {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok && c.CorrelationID != nil {
		return *c.CorrelationID, true
	}

	return "", false
}

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok {
		c.UserID = &userID
		return context.WithValue(ctx, ReqCtxKey, c)
	}

	return context.WithValue(ctx, ReqCtxKey, ReqCtx{UserID: &userID})
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok && c.UserID != nil {
		return *c.UserID, true
	}

	return uuid.UUID{}, false
}

func WithResponseCode(ctx context.Context, responseCode int) context.Context {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok {
		c.ResponseCode = &responseCode
		return context.WithValue(ctx, ReqCtxKey, c)
	}

	return context.WithValue(ctx, ReqCtxKey, ReqCtx{ResponseCode: &responseCode})
}

func GetResponseCode(ctx context.Context) (int, bool) {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok && c.ResponseCode != nil {
		return *c.ResponseCode, true
	}

	return 0, false
}

func WithClaimID(ctx context.Context, claimID uint64) context.Context {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok {
		c.ClaimID = &claimID
		return context.WithValue(ctx, ReqCtxKey, c)
	}

	return context.WithValue(ctx, ReqCtxKey, ReqCtx{ClaimID: &claimID})
}

func GetClaimID(ctx context.Context) (uint64, bool) {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok && c.ClaimID != nil {
		return *c.ClaimID, true
	}

	return 0, false
}

func WithRole(ctx context.Context, r roles.Role) context.Context {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok {
		c.Role = &r
		return context.WithValue(ctx, ReqCtxKey, c)
	}

	return context.WithValue(ctx, ReqCtxKey, ReqCtx{Role: &r})
}

func GetRole(ctx context.Context) (roles.Role, bool) {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok && c.Role != nil {
		return *c.Role, true
	}

	return roles.Unknown, false
}
