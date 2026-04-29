package handlers

import (
	"at-backend-claims/internal/pkg/apperror"
	"at-backend-claims/internal/pkg/logs"
	"at-backend-claims/internal/pkg/responses"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (h claimsHandlerForAdmins) getClaimByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.NotFound(r.Context(), w)
		return
	}

	claim, err := h.usecase.GetByID(r.Context(), id)
	if err != nil {
		logs.Error(r.Context(), err)

		switch {
		case errors.Is(err, apperror.ErrClaimNotExists):
			responses.NotFound(r.Context(), w)
		default:
			responses.InternalServerError(r.Context(), w)
		}

		return
	}

	resp := map[string]any{
		"id":                  id,
		"title":               claim.Title,
		"description":         claim.Description,
		"category":            claim.Category,
		"status":              claim.Status,
		"feedback":            claim.Feedback,
		"photos":              claim.Photos,
		"latitude":            claim.Latitude,
		"longitude":           claim.Longitude,
		"created_at":          claim.CreatedAt,
		"updated_at":          claim.UpdatedAt,
		"status_updated_at":   claim.StatusUpdatedAt,
		"feedback_updated_at": claim.FeedbackUpdatedAt,
		"created_by":          claim.CreatedBy,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}

func (h claimsHandlerForAdmins) getPageClaims(w http.ResponseWriter, r *http.Request) {
	pageSize := getPageSize(r.URL.Query())

	cursorString := r.URL.Query().Get("cursor")
	if cursorString == "" {
		h.getFirstPageClaims(w, r, pageSize)
		return
	}

	cursor, err := strconv.ParseUint(cursorString, 10, 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	page, err := h.usecase.GetPage(r.Context(), cursor, pageSize)
	if err != nil {
		logs.Error(r.Context(), err)

		responses.InternalServerError(r.Context(), w)
		return
	}

	pageResp := sliceClaimsToSliceClaimsForPage(page)

	resp := map[string]any{
		"claims": pageResp,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}

func (h claimsHandlerForAdmins) getFirstPageClaims(w http.ResponseWriter, r *http.Request, pageSize uint64) {
	page, err := h.usecase.GetFirstPage(r.Context(), pageSize)
	if err != nil {
		logs.Error(r.Context(), err)

		responses.InternalServerError(r.Context(), w)
		return
	}

	pageResp := sliceClaimsToSliceClaimsForPage(page)

	resp := map[string]any{
		"claims": pageResp,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}

func (h claimsHandlerForAdmins) getUserClaimsPage(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r.URL.Query())
	if !ok {
		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	pageSize := getPageSize(r.URL.Query())

	cursorString := r.URL.Query().Get("cursor")
	if cursorString == "" {
		h.getUserClaimsFirstPage(w, r, pageSize, userID)
		return
	}

	cursor, err := strconv.ParseUint(cursorString, 10, 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	page, err := h.usecase.GetUserPage(r.Context(), cursor, userID, pageSize)
	if err != nil {
		logs.Error(r.Context(), err)

		responses.InternalServerError(r.Context(), w)
		return
	}

	pageResp := sliceClaimsToSliceClaimsForPage(page)

	resp := map[string]any{
		"claims": pageResp,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}

func (h claimsHandlerForAdmins) getUserClaimsFirstPage(w http.ResponseWriter, r *http.Request, pageSize uint64, uid uuid.UUID) {
	page, err := h.usecase.GetFirstUserPage(r.Context(), uid, pageSize)
	if err != nil {
		logs.Error(r.Context(), err)

		responses.InternalServerError(r.Context(), w)
		return
	}

	pageResp := sliceClaimsToSliceClaimsForPage(page)

	resp := map[string]any{
		"claims": pageResp,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}
