package usecase

import "context"

func (c *categoriesUsecase) CreateCategory(ctx context.Context, name string) (uint64, error) {
	id, err := c.repo.CreateCategory(ctx, name)
	if err != nil {
		return 0, err
	}

	c.unloadCategories()

	return id, nil
}

func (c *categoriesUsecase) CreateSubcategory(ctx context.Context, categoryID uint64, name string) (uint64, error) {
	id, err := c.repo.CreateSubcategory(ctx, categoryID, name)
	if err != nil {
		return 0, err
	}

	c.unloadCategories()

	return id, nil
}
