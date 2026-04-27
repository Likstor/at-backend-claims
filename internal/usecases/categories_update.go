package usecase

import "context"

func (c *categoriesUsecase) UpdateCategory(ctx context.Context, categoryID uint64, name string) error {
	if err := c.repo.UpdateCategory(ctx, categoryID, name); err != nil {
		return err
	}
	
	c.unloadCategories()
	
	return nil
}

func (c *categoriesUsecase) UpdateSubcategory(ctx context.Context, subcategoryID uint64, name string) error {
	if err :=  c.repo.UpdateSubcategory(ctx, subcategoryID, name); err != nil {
		return err
	}

	c.unloadCategories()

	return nil
}
