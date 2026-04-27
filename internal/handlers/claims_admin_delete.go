package handlers

import (
	"at-backend-claims/internal/pkg/apperror"
	"at-backend-claims/internal/pkg/logs"
	"at-backend-claims/internal/pkg/responses"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

func (h claimsHandlerForAdmins) deleteClaim(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.NotFound(r.Context(), w)
		return
	}

	if err := h.usecase.Delete(r.Context(), id); err != nil {
		logs.Error(r.Context(), err)

		switch {
		case errors.Is(err, apperror.ErrClaimNotExists):
			responses.NotFound(r.Context(), w)
		default:
			responses.InternalServerError(r.Context(), w)
		}

		return

	}

	w.WriteHeader(http.StatusNoContent)
}
