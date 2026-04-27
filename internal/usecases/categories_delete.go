package usecase

import "context"

func (c *categoriesUsecase) DeleteCategory(ctx context.Context, categoryID uint64) error {
	if err := c.repo.DeleteCategory(ctx, categoryID); err != nil {
		return err
	}
	
	c.unloadCategories()
	
	return nil
}

func (c *categoriesUsecase) DeleteSubcategory(ctx context.Context, subcategoryID uint64) error {
	if err := c.repo.DeleteSubcategory(ctx, subcategoryID); err != nil {
		return err
	}
	
	c.unloadCategories()
	
	return nil
}
