package usecase

import (
	"at-backend-claims/internal/domain"
	"context"
)

func (c *categoriesUsecase) GetAll(ctx context.Context) ([]domain.Category, error) {
	c.mutex.RLock()
	if c.categories != nil {
		defer c.mutex.RUnlock()
		return c.categories, nil
	}
	c.mutex.RUnlock()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	res, err := c.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	catMap := make(map[string]struct{})
	subcatMap := make(map[string]struct{})

	for _, cat := range res {
		catMap[cat.Name] = struct{}{}

		for _, subcat := range cat.Subcategories {
			subcatMap[subcat.Name] = struct{}{}
		}
	}

	c.categoriesMap = catMap
	c.subcategoriesMap = subcatMap
	c.categories = res

	return res, nil
}
