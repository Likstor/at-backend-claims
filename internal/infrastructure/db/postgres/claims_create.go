package postgres

import (
	"at-backend-claims/internal/domain"
	"at-backend-claims/internal/pkg/apperror"
	"context"
)

const queryClaimsCreate = `
	INSERT INTO claims 
		(created_by, title, description, category, status, photos, latitude, longitude) 
	VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8) 
	RETURNING id
`

func (cr claimsRepository) Create(ctx context.Context, data domain.Claim) (uint64, error) {
	var id uint64

	if err := cr.Conn(ctx).QueryRow(
		ctx,
		queryClaimsCreate,
		data.CreatedBy,
		data.Title,
		data.Description,
		data.Category,
		data.Status.String(),
		data.Photos,
		data.Latitude,
		data.Longitude,
	).Scan(&id); err != nil {
		return id, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	return id, nil
}
