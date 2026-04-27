package postgres

import (
	"at-backend-claims/internal/domain"
	"at-backend-claims/internal/pkg/apperror"
	"context"
	"log/slog"
)

const queryCategoriesGetAll = `
	SELECT
		c.id, c.name, s.id, s.name
	FROM categories AS c
	LEFT JOIN subcategories AS s
		ON s.category_id = c.id
	ORDER BY c.id, s.id
`

func (r categoriesRepository) GetAll(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.Conn(ctx).Query(ctx, queryCategoriesGetAll)
	if err != nil {
		return nil, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}
	defer rows.Close()

	result := make([]domain.Category, 0)
	subResult := make([]domain.Subcategory, 0)

	var curr uint64 = 0
	for rows.Next() {
		var (
			categoryID      uint64
			categoryName    string
			subcategoryID   *uint64
			subcategoryName *string
		)

		if err := rows.Scan(&categoryID, &categoryName, &subcategoryID, &subcategoryName); err != nil {
			slog.WarnContext(ctx, err.Error())
			continue
		}

		if curr != categoryID {
			if len(result) > 0 {
				result[len(result)-1].Subcategories = subResult
			}

			result = append(result, domain.Category{
				ID:            categoryID,
				Name:          categoryName,
				Subcategories: make([]domain.Subcategory, 0),
			})

			subResult = make([]domain.Subcategory, 0)
			curr = categoryID
		}

		if subcategoryID != nil && subcategoryName != nil && *subcategoryName != "" {
			subResult = append(subResult, domain.Subcategory{
				ID:   *subcategoryID,
				Name: *subcategoryName,
			})
		}
	}

	if len(result) > 0 {
		result[len(result)-1].Subcategories = subResult
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	return result, nil

}
