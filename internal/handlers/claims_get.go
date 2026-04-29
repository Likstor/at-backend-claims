package handlers

import (
	"at-backend-claims/internal/pkg/apperror"
	"at-backend-claims/internal/pkg/logs"
	"at-backend-claims/internal/pkg/reqctx"
	"at-backend-claims/internal/pkg/responses"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

func (h claimsHandler) getById(w http.ResponseWriter, r *http.Request) {
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
		"id":                  claim.ID,
		"title":               claim.Title,
		"description":         claim.Description,
		"category":            claim.Category,
		"status":              claim.Status,
		"photos":              claim.Photos,
		"latitude":            claim.Latitude,
		"longitude":           claim.Longitude,
		"feedback":            claim.Feedback,
		"created_at":          claim.CreatedAt,
		"updated_at":          claim.UpdatedAt,
		"status_updated_at":   claim.StatusUpdatedAt,
		"feedback_updated_at": claim.FeedbackUpdatedAt,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}

func (h claimsHandler) getPage(w http.ResponseWriter, r *http.Request) {
	pageSize := getPageSize(r.URL.Query())

	cursorString := r.URL.Query().Get("cursor")
	if cursorString == "" {
		h.getFirstPage(w, r, pageSize)
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

	uid, ok := reqctx.GetUserID(r.Context())
	if !ok {
		slog.ErrorContext(r.Context(), apperror.ErrCtxEmptyUserID.Error())

		responses.InternalServerError(r.Context(), w)
		return
	}

	page, err := h.usecase.GetUserPage(r.Context(), cursor, uid, pageSize)
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

func (h claimsHandler) getFirstPage(w http.ResponseWriter, r *http.Request, pageSize uint64) {
	uid, ok := reqctx.GetUserID(r.Context())
	if !ok {
		slog.ErrorContext(r.Context(), apperror.ErrCtxEmptyUserID.Error())

		responses.InternalServerError(r.Context(), w)
		return
	}

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

func (c claimsHandler) getClaimsByArea(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	lat1 := query.Get("lat1")
	long1 := query.Get("long1")
	lat2 := query.Get("lat2")
	long2 := query.Get("long2")

	if lat1 == "" || long1 == "" || lat2 == "" || long2 == "" {
		slog.WarnContext(r.Context(), "empty query parameters")

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	lat1float, err := strconv.ParseFloat(lat1, 64)
	if err != nil {
		slog.WarnContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	long1float, err := strconv.ParseFloat(long1, 64)
	if err != nil {
		slog.WarnContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	lat2float, err := strconv.ParseFloat(lat2, 64)
	if err != nil {
		slog.WarnContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	long2float, err := strconv.ParseFloat(long2, 64)
	if err != nil {
		slog.WarnContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	claims, err := c.usecase.GetByArea(r.Context(), lat1float, long1float, lat2float, long2float)
	if err != nil {
		responses.InternalServerError(r.Context(), w)
		return
	}

	claimsForMap := make([]map[string]any, len(claims))

	for _, claim := range claims {
		claimsForMap = append(claimsForMap, map[string]any{
			"id":        claim.ID,
			"latitude":  claim.Latitude,
			"longitude": claim.Longitude,
		})
	}

	resp := map[string]any{
		"claims": claimsForMap,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}
