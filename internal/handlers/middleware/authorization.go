package middleware

import (
	"at-backend-claims/internal/pkg/reqctx"
	"at-backend-claims/internal/pkg/responses"
	"at-backend-claims/internal/pkg/roles"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Verifier func(context context.Context, token string) (uuid.UUID, roles.Role, error)

func Authorization(next http.Handler, minRole roles.Role, verifier Verifier) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			slog.WarnContext(r.Context(), "authorization header empty")

			responses.Error(
				r.Context(),
				w,
				http.StatusBadRequest,
				"authorization header empty",
			)
			return
		}

		authHeaderParts := strings.Fields(authHeader)
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			slog.WarnContext(r.Context(), "invalid value in authorization header")

			responses.Error(
				r.Context(),
				w,
				http.StatusBadRequest,
				"400 authorization header must be \"Bearer {token}\"",
			)
			return
		}

		userID, role, err := verifier(r.Context(), authHeaderParts[1])
		if err != nil {
			slog.WarnContext(r.Context(), err.Error())

			responses.Error(
				r.Context(),
				w,
				http.StatusUnauthorized,
				http.StatusText(http.StatusUnauthorized),
			)
			return
		}

		if roles.IsRoleALowerThanB(role, minRole) {
			slog.WarnContext(r.Context(), fmt.Sprintf("role %s < %s", role.String(), minRole.String()))

			responses.Error(
				r.Context(),
				w,
				http.StatusForbidden,
				http.StatusText(http.StatusForbidden),
			)
			return
		}

		ctx := reqctx.WithUserID(r.Context(), userID)
		ctx = reqctx.WithRole(ctx, role)

		r = r.WithContext(ctx)

		slog.InfoContext(r.Context(), "authorized success")

		next.ServeHTTP(w, r)
	})
}
