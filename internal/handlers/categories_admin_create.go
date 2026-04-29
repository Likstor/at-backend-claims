package handlers

import (
	"at-backend-claims/internal/pkg/apperror"
	"at-backend-claims/internal/pkg/logs"
	"at-backend-claims/internal/pkg/responses"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type categoryCreateRequest struct {
	Name *string `json:"category_name"`
}

func (s categoryCreateRequest) Validate() error {
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

func (h categoriesHandlerForAdmins) createCategory(w http.ResponseWriter, r *http.Request) {
	var categoryDTO categoryCreateRequest

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

	categoryID, err := h.usecase.CreateCategory(r.Context(), *categoryDTO.Name)
	if err != nil {
		logs.Error(r.Context(), err)

		switch {
		case errors.Is(err, apperror.ErrCategoryAlreadyExists):
			responses.Error(
				r.Context(),
				w,
				http.StatusConflict,
				err.Error(),
			)
		default:
			responses.InternalServerError(r.Context(), w)
		}

		return
	}

	resp := map[string]any{
		"category_id": categoryID,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}

type subcategoryCreateRequest struct {
	CategoryID *uint64 `json:"category_id"`
	Name       *string `json:"subcategory_name"`
}

func (s subcategoryCreateRequest) Validate() error {
	var err error

	if s.Name == nil {
		err = errors.Join(err, errors.New("subcategory_name field is missing"))
	} else {
		if *s.Name == "" {
			err = errors.Join(err, errors.New("subcategory_name field is empty"))
		}
	}

	if s.CategoryID == nil {
		err = errors.Join(err, errors.New("category_id field is missing"))
	}

	return err
}

func (h categoriesHandlerForAdmins) createSubcategory(w http.ResponseWriter, r *http.Request) {
	var subcategoryDTO subcategoryCreateRequest

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

	categoryID, err := h.usecase.CreateSubcategory(r.Context(), *subcategoryDTO.CategoryID, *subcategoryDTO.Name)
	if err != nil {
		logs.Error(r.Context(), err)

		switch {
		case errors.Is(err, apperror.ErrSubcategoryAlreadyExists):
			responses.Error(
				r.Context(),
				w,
				http.StatusConflict,
				err.Error(),
			)
		case errors.Is(err, apperror.ErrCategoryNotExists):
			responses.NotFound(r.Context(), w)
		default:
			responses.InternalServerError(r.Context(), w)
		}

		return
	}

	resp := map[string]any{
		"subcategory_id": categoryID,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}
