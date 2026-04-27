package handlers

import (
	"at-backend-claims/internal/domain"
	"at-backend-claims/internal/pkg/apperror"
	"at-backend-claims/internal/pkg/logs"
	"at-backend-claims/internal/pkg/reqctx"
	"at-backend-claims/internal/pkg/responses"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

type updateClaimRequest struct {
	claimWithoutCategory
}

func (claim updateClaimRequest) Validate() error {
	var err error

	err = errors.Join(err, claim.claimWithoutCategory.Validate())

	return err
}

// TODO: Необходимо уточнить поля, которые можно будет менять

func (h claimsHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.NotFound(r.Context(), w)
		return
	}

	var claimDTO updateClaimRequest

	if err = json.NewDecoder(r.Body).Decode(&claimDTO); err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	if err = claimDTO.Validate(); err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusUnprocessableEntity,
			err.Error(),
		)
		return
	}

	uid, ok := reqctx.GetUserID(r.Context())
	if !ok {
		slog.ErrorContext(r.Context(), "empty userID in context")

		responses.InternalServerError(r.Context(), w)
		return
	}

	claim := domain.Claim{
		ID:          id,
		CreatedBy:   uid,
		Title:       *claimDTO.Title,
		Description: *claimDTO.Description,
		Latitude:    *claimDTO.Latitude,
		Longitude:   *claimDTO.Longitude,
	}

	if err := h.usecase.Update(r.Context(), claim); err != nil {
		logs.Error(r.Context(), err)

		switch {
		case errors.Is(err, apperror.ErrUnknownClaimCategory):
			responses.Error(
				r.Context(),
				w,
				http.StatusUnprocessableEntity,
				err.Error(),
			)
		case errors.Is(err, apperror.ErrUserCannotPerformOperation):
			responses.Error(
				r.Context(),
				w,
				http.StatusForbidden,
				http.StatusText(http.StatusForbidden),
			)
		case errors.Is(err, apperror.ErrClaimNotExists):
			responses.NotFound(r.Context(), w)
		default:
			responses.InternalServerError(r.Context(), w)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
