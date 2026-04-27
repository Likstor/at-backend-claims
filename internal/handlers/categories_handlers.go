package handlers

import (
	"at-backend-claims/internal/domain"
	"at-backend-claims/internal/handlers/middleware"
	"at-backend-claims/internal/pkg/roles"
	"context"
	"fmt"
	"net/http"
)

type categoriesUsecase interface {
	GetAll(ctx context.Context) ([]domain.Category, error)
}

type categoriesHandler struct {
	usecase categoriesUsecase
}

func NewCategoriesHandler(uc categoriesUsecase) *categoriesHandler {
	return &categoriesHandler{
		usecase: uc,
	}
}

func (c categoriesHandler) Setup(prefix string, verifier middleware.Verifier, mux *http.ServeMux) {
	muxWithAuth := http.NewServeMux()

	muxWithAuth.HandleFunc("GET /", c.getAll)

	muxWithAuthWrapped := middleware.Authorization(muxWithAuth, roles.User, verifier)

	mux.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(prefix, muxWithAuthWrapped))
}

type categoriesUsecaseForAdmins interface {
	categoriesUsecase

	CreateCategory(ctx context.Context, name string) (uint64, error)
	CreateSubcategory(ctx context.Context, categoryID uint64, name string) (uint64, error)

	UpdateCategory(ctx context.Context, categoryID uint64, name string) error
	UpdateSubcategory(ctx context.Context, subcategoryID uint64, name string) error

	DeleteCategory(ctx context.Context, categoryID uint64) error
	DeleteSubcategory(ctx context.Context, subcategoryID uint64) error
}

type categoriesHandlerForAdmins struct {
	usecase categoriesUsecaseForAdmins
}

func NewCategoriesHandlerForAdmins(uc categoriesUsecaseForAdmins) *categoriesHandlerForAdmins {
	return &categoriesHandlerForAdmins{
		usecase: uc,
	}
}

func (c categoriesHandlerForAdmins) Setup(prefix string, verifier middleware.Verifier, mux *http.ServeMux) {
	muxWithAuth := http.NewServeMux()

	muxWithAuth.HandleFunc("POST /", c.createCategory)
	muxWithAuth.HandleFunc("PUT /{id}", c.updateCategory)
	muxWithAuth.HandleFunc("DELETE /{id}", c.deleteCategory)

	muxWithAuth.HandleFunc("POST /subcategories/", c.createSubcategory)
	muxWithAuth.HandleFunc("PUT /subcategories/{id}", c.updateSubcategory)
	muxWithAuth.HandleFunc("DELETE /subcategories/{id}", c.deleteSubcategory)

	muxWithAuthWrapped := middleware.Authorization(muxWithAuth, roles.Admin, verifier)

	mux.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(prefix, muxWithAuthWrapped))
}