package postgres

import (
	"at-backend-claims/internal/pkg/apperror"
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

const queryCategoriesCreateCategory = `
	INSERT INTO categories
		(name)
	VALUES
		($1)
	RETURNING id
`

func (r categoriesRepository) CreateCategory(ctx context.Context, name string) (uint64, error) {
	var id uint64

	if err := r.Conn(ctx).QueryRow(ctx, queryCategoriesCreateCategory, name).Scan(&id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return 0, apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrCategoryAlreadyExists)
			}
		}

		return 0, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	return id, nil
}

const queryCategoriesCreateSubcategory = `
	INSERT INTO subcategories
		(name, category_id)
	VALUES
		($1, $2)
	RETURNING id
`

func (r categoriesRepository) CreateSubcategory(ctx context.Context, categoryID uint64, name string) (uint64, error) {
	var id uint64

	if err := r.Conn(ctx).QueryRow(ctx, queryCategoriesCreateSubcategory, name, categoryID).Scan(&id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return 0, apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrSubcategoryAlreadyExists)
			case pgerrcode.ForeignKeyViolation:
				return 0, apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrCategoryNotExists)
			}

			return 0, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
		}

		return 0, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	return id, nil
}
