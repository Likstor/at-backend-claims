package usecase

import "context"

func (c *categoriesUsecase) CategoryExist(ctx context.Context, category string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if _, ok := c.categoriesMap[category]; ok {
		return true
	}

	return false
}

func (c *categoriesUsecase) SubcategoryExist(ctx context.Context, subcategory string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if _, ok := c.subcategoriesMap[subcategory]; ok {
		return true
	}

	return false
}
