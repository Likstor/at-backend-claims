package postgres

import (
	"at-backend-claims/internal/pkg/apperror"
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

const queryCategoriesUpdateCategory = `
	UPDATE categories
	SET
		name = $2
	WHERE
		id = $1
`

func (r categoriesRepository) UpdateCategory(ctx context.Context, categoryID uint64, name string) error {
	resp, err := r.Conn(ctx).Exec(ctx, queryCategoriesUpdateCategory, categoryID, name)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrCategoryAlreadyExists)
			}
		}

		return apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	if resp.RowsAffected() < 1 {
		return apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrCategoryNotExists)
	}

	return nil
}

const queryCategoriesUpdateSubcategory = `
	UPDATE subcategories
	SET 
		name = $2
	WHERE
		id = $1
`

func (r categoriesRepository) UpdateSubcategory(ctx context.Context, subcategoryID uint64, name string) error {
	resp, err := r.Conn(ctx).Exec(ctx, queryCategoriesUpdateSubcategory, subcategoryID, name)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrSubcategoryAlreadyExists)
			}
		}

		return apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	if resp.RowsAffected() < 1 {
		return apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrCategoryNotExists)
	}

	return nil
}
