package usecase

import (
	"at-backend-claims/internal/domain"
	"context"
	"sync"
)

var (
	dummyCategory    = domain.Category{}
	dummySubcategory = domain.Subcategory{}
)

type categoriesRepo interface {
	GetAll(ctx context.Context) ([]domain.Category, error)

	CreateCategory(ctx context.Context, name string) (uint64, error)
	CreateSubcategory(ctx context.Context, categoryID uint64, name string) (uint64, error)

	DeleteCategory(ctx context.Context, categoryID uint64) error
	DeleteSubcategory(ctx context.Context, subcategoryID uint64) error

	UpdateCategory(ctx context.Context, categoryID uint64, name string) error
	UpdateSubcategory(ctx context.Context, subcategoryID uint64, name string) error
}

type categoriesUsecase struct {
	repo        categoriesRepo
	maxPageSize uint64

	categories []domain.Category
	mutex      sync.RWMutex

	categoriesMap    map[string]struct{}
	subcategoriesMap map[string]struct{}
}

func NewCategoriesUsecase(ctx context.Context, repo categoriesRepo, maxPageSize uint64) *categoriesUsecase {
	cu := &categoriesUsecase{
		repo:        repo,
		maxPageSize: maxPageSize,

		mutex: sync.RWMutex{},
	}

	cu.GetAll(ctx) // Инициализация categoriesMap, subcategoriesMap

	return cu
}

func (c *categoriesUsecase) unloadCategories() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.categories = nil
}
