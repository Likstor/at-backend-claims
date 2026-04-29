package handlers

import (
	"at-backend-claims/internal/domain"
	"at-backend-claims/internal/pkg/apperror"
	"at-backend-claims/internal/pkg/logs"
	"at-backend-claims/internal/pkg/reqctx"
	"at-backend-claims/internal/pkg/responses"
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

var defaultMaxMemory int64 = 10 << 20

type createProblemRequest struct {
	claimWithoutCategory
	Category *string `json:"category"`
}

func (claim createProblemRequest) Validate() error {
	var err error

	err = errors.Join(err, claim.claimWithoutCategory.Validate())

	if claim.Category == nil {
		err = errors.Join(err, errors.New("category field is missing"))
	} else {
		if *claim.Category == "" {
			err = errors.Join(err, errors.New("category files is empty"))
		}
	}

	return err
}

func (h claimsHandler) createProblem(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(int64(defaultMaxMemory)); err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusRequestEntityTooLarge,
			http.StatusText(http.StatusRequestEntityTooLarge),
		)
		return
	}

	lat, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusUnprocessableEntity,
			http.StatusText(http.StatusUnprocessableEntity),
		)
		return
	}

	long, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusUnprocessableEntity,
			http.StatusText(http.StatusUnprocessableEntity),
		)
		return
	}

	var claimDTO = createProblemRequest{
		claimWithoutCategory: claimWithoutCategory{
			Title:       new(r.FormValue("title")),
			Description: new(r.FormValue("description")),
			Longitude:   &long,
			Latitude:    &lat,
		},
		Category: new(r.FormValue("category")),
	}

	if err := claimDTO.Validate(); err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusUnprocessableEntity,
			http.StatusText(http.StatusUnprocessableEntity),
		)
		return
	}

	uid, ok := reqctx.GetUserID(r.Context())
	if !ok {
		slog.ErrorContext(r.Context(), apperror.ErrCtxEmptyUserID.Error())

		responses.InternalServerError(r.Context(), w)
		return
	}

	claim := domain.Claim{
		CreatedBy:   uid,
		Title:       *claimDTO.Title,
		Description: *claimDTO.Description,
		Category:    *claimDTO.Category,
		Photos:      nil,
		Latitude:    *claimDTO.Latitude,
		Longitude:   *claimDTO.Longitude,
	}

	files := make([][]byte, 0, len(r.MultipartForm.File))

	for fn := range r.MultipartForm.File {
		file, _, err := r.FormFile(fn)
		if err != nil {
			slog.ErrorContext(r.Context(), err.Error())
			continue
		}
		defer file.Close()

		buf := &bytes.Buffer{}

		if _, err := io.Copy(buf, file); err != nil {
			slog.ErrorContext(r.Context(), err.Error())
			continue
		}

		files = append(files, buf.Bytes())
	}

	id, err := h.usecase.CreateProblem(r.Context(), claim, files)
	if err != nil {
		logs.Error(r.Context(), err)

		switch {
		case errors.Is(err, apperror.ErrUnknownClaimCategory), errors.Is(err, apperror.ErrPointIsNotInPolygon):
			responses.Error(
				r.Context(),
				w,
				http.StatusUnprocessableEntity,
				http.StatusText(http.StatusUnprocessableEntity),
			)
		default:
			responses.InternalServerError(r.Context(), w)
		}

		return
	}

	resp := map[string]any{
		"claim_id": id,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}

type createProposalRequest struct {
	claimWithoutCategory
}

func (claim createProposalRequest) Validate() error {
	var err error

	err = errors.Join(err, claim.claimWithoutCategory.Validate())

	return err
}

func (h claimsHandler) createProposal(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(int64(defaultMaxMemory)); err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusRequestEntityTooLarge,
			http.StatusText(http.StatusRequestEntityTooLarge),
		)
		return
	}

	lat, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusUnprocessableEntity,
			http.StatusText(http.StatusUnprocessableEntity),
		)
		return
	}

	long, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusUnprocessableEntity,
			http.StatusText(http.StatusUnprocessableEntity),
		)
		return
	}

	var claimDTO = createProposalRequest{
		claimWithoutCategory: claimWithoutCategory{
			Title:       new(r.FormValue("title")),
			Description: new(r.FormValue("description")),
			Longitude:   &long,
			Latitude:    &lat,
		},
	}

	if err := claimDTO.Validate(); err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusUnprocessableEntity,
			http.StatusText(http.StatusUnprocessableEntity),
		)
		return
	}

	uid, ok := reqctx.GetUserID(r.Context())
	if !ok {
		slog.ErrorContext(r.Context(), apperror.ErrCtxEmptyUserID.Error())

		responses.InternalServerError(r.Context(), w)
		return
	}

	claim := domain.Claim{
		CreatedBy:   uid,
		Title:       *claimDTO.Title,
		Description: *claimDTO.Description,
		Latitude:    *claimDTO.Latitude,
		Longitude:   *claimDTO.Longitude,
	}

	id, err := h.usecase.CreateProposal(r.Context(), claim)
	if err != nil {
		logs.Error(r.Context(), err)

		switch {
		case errors.Is(err, apperror.ErrPointIsNotInPolygon):
			responses.Error(
				r.Context(),
				w,
				http.StatusUnprocessableEntity,
				http.StatusText(http.StatusUnprocessableEntity),
			)
		default:
			responses.InternalServerError(r.Context(), w)
		}

		return
	}

	resp := map[string]any{
		"claim_id": id,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}
