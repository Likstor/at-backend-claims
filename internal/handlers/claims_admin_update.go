package handlers

import (
	"at-backend-claims/internal/domain"
	"at-backend-claims/internal/pkg/apperror"
	"at-backend-claims/internal/pkg/logs"
	"at-backend-claims/internal/pkg/responses"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

type claimStatusUpdateRequest struct {
	Status *domain.ClaimStatus `json:"status"`
}

func (claim claimStatusUpdateRequest) Validate() error {
	var err error

	if claim.Status == nil {
		err = errors.Join(err, errors.New("status field is missing or null"))
	} else {
		if *claim.Status == "" {
			err = errors.Join(err, errors.New("status field is empty"))
		}
	}

	return err
}

func (h claimsHandlerForAdmins) changeClaimStatus(w http.ResponseWriter, r *http.Request) {
	var statusDTO claimStatusUpdateRequest

	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.NotFound(r.Context(), w)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&statusDTO); err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)

		return
	}

	if err := statusDTO.Validate(); err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusUnprocessableEntity,
			err.Error(),
		)

		return
	}

	if err := h.usecase.ChangeStatus(r.Context(), id, *statusDTO.Status); err != nil {
		logs.Error(r.Context(), err)

		switch {
		case errors.Is(err, apperror.ErrUnknownClaimStatus):
			responses.Error(
				r.Context(),
				w,
				http.StatusUnprocessableEntity,
				err.Error(),
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

type claimFeedbackRequest struct {
	Feedback *string `json:"feedback"`
}

func (claim claimFeedbackRequest) Validate() error {
	var err error

	if claim.Feedback == nil {
		err = errors.Join(err, errors.New("feedback field is missing or null"))
	} else {
		if *claim.Feedback == "" {
			err = errors.Join(err, errors.New("feedback field is empty"))
		}
	}

	return err
}

func (h claimsHandlerForAdmins) addFeedbackToClaim(w http.ResponseWriter, r *http.Request) {
	var feedbackDTO claimFeedbackRequest

	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.NotFound(r.Context(), w)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&feedbackDTO); err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)

		return
	}

	if err := feedbackDTO.Validate(); err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusUnprocessableEntity,
			err.Error(),
		)

		return
	}

	if err := h.usecase.AddFeedback(r.Context(), id, *feedbackDTO.Feedback); err != nil {
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
