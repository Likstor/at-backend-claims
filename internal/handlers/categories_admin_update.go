package handlers

import (
	"at-backend-claims/internal/pkg/apperror"
	"at-backend-claims/internal/pkg/logs"
	"at-backend-claims/internal/pkg/responses"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

type categoryUpdateRequest struct {
	Name *string `json:"category_name"`
}

func (s categoryUpdateRequest) Validate() error {
	var err error

	if s.Name == nil {
		err = errors.Join(err, errors.New("category_name field is missing"))
	} else {
		if *s.Name == "" {
			err = errors.Join(err, errors.New("category_name field is empty"))
		}
	}

	return err
}

func (h categoriesHandlerForAdmins) updateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.NotFound(r.Context(), w)
		return
	}

	var categoryDTO categoryUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&categoryDTO); err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	if err := categoryDTO.Validate(); err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusUnprocessableEntity,
			http.StatusText(http.StatusUnprocessableEntity),
		)
		return
	}

	if err := h.usecase.UpdateSubcategory(r.Context(), id, *categoryDTO.Name); err != nil {
		logs.Error(r.Context(), err)

		switch {
		case errors.Is(err, apperror.ErrCategoryNotExists):
			responses.NotFound(r.Context(), w)
		case errors.Is(err, apperror.ErrCategoryAlreadyExists):
			responses.Error(
				r.Context(),
				w,
				http.StatusConflict,
				apperror.ErrCategoryAlreadyExists.Error(),
			)
		default:
			responses.InternalServerError(r.Context(), w)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type subcategoryUpdateRequest struct {
	Name *string `json:"subcategory_name"`
}

func (s subcategoryUpdateRequest) Validate() error {
	var err error

	if s.Name == nil {
		err = errors.Join(err, errors.New("subcategory_name field is missing"))
	} else {
		if *s.Name == "" {
			err = errors.Join(err, errors.New("subcategory_name field is empty"))
		}
	}

	return err
}

func (h categoriesHandlerForAdmins) updateSubcategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.NotFound(r.Context(), w)
		return
	}

	var subcategoryDTO subcategoryUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&subcategoryDTO); err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	if err := subcategoryDTO.Validate(); err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusUnprocessableEntity,
			http.StatusText(http.StatusUnprocessableEntity),
		)
		return
	}

	if err := h.usecase.UpdateSubcategory(r.Context(), id, *subcategoryDTO.Name); err != nil {
		logs.Error(r.Context(), err)

		switch {
		case errors.Is(err, apperror.ErrCategoryNotExists):
			responses.NotFound(r.Context(), w)
		case errors.Is(err, apperror.ErrSubcategoryAlreadyExists):
			responses.Error(
				r.Context(),
				w,
				http.StatusConflict,
				apperror.ErrSubcategoryAlreadyExists.Error(),
			)
		default:
			responses.InternalServerError(r.Context(), w)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
