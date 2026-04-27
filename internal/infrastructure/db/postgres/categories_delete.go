package postgres

import (
	"at-backend-claims/internal/pkg/apperror"
	"context"
)

func (r categoriesRepository) delete(ctx context.Context, query string, id uint64) error {
	resp, err := r.Conn(ctx).Exec(ctx, query, id)
	if err != nil {
		return apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	if resp.RowsAffected() < 1 {
		return apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrCategoryNotExists)
	}

	return nil
}

const queryCategoriesDeleteCategory = `
	DELETE FROM categories
	WHERE id = $1
`

func (r categoriesRepository) DeleteCategory(ctx context.Context, categoryID uint64) error {
	if err := r.delete(ctx, queryCategoriesDeleteCategory, categoryID); err != nil {
		return err
	}

	return nil
}

const queryCategoriesDeleteSubcategory = `
	DELETE FROM subcategories
	WHERE id = $1
`

func (r categoriesRepository) DeleteSubcategory(ctx context.Context, subcategoryID uint64) error {
	if err := r.delete(ctx, queryCategoriesDeleteSubcategory, subcategoryID); err != nil {
		return err
	}

	return nil
}
