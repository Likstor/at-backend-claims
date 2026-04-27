package postgres

import (
	"at-backend-claims/internal/domain"
	"at-backend-claims/internal/pkg/apperror"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const queryClaimsGetByID = `
	SELECT 
		title,
		created_by,
		description, 
		category, 
		status,
		photos, 
		latitude, 
		longitude, 
		feedback, 
		created_at, 
		updated_at, 
		status_updated_at, 
		feedback_updated_at
	FROM claims
	WHERE id = $1
`

func (cr claimsRepository) GetByID(ctx context.Context, id uint64) (domain.Claim, error) {
	rows, err := cr.Conn(ctx).Query(ctx, queryClaimsGetByID, id)
	if err != nil {
		return dummyClaim, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}
	defer rows.Close()

	claim, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[domain.Claim])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dummyClaim, apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrClaimNotExists)
		}

		return dummyClaim, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	claim.ID = id

	return claim, nil
}

func (cr claimsRepository) getClaimsFromRows(ctx context.Context, rows pgx.Rows) ([]domain.Claim, error) {
	claims, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.Claim])
	if err != nil {
		return nil, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}
	defer rows.Close()

	return claims, nil
}

const queryClaimsGetFirstPage = `
	SELECT 
		id, title, description, category, status, latitude, longitude, created_at, updated_at
	FROM claims
	ORDER BY id DESC
	LIMIT $1
`

func (cr claimsRepository) GetFirstPage(ctx context.Context, pageSize uint64) ([]domain.Claim, error) {
	rows, err := cr.Conn(ctx).Query(ctx, queryClaimsGetFirstPage, pageSize)
	if err != nil {
		return nil, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}
	defer rows.Close()

	return cr.getClaimsFromRows(ctx, rows)
}

const queryClaimsGetPage = `
	SELECT 
		id, title, description, category, status, latitude, longitude, created_at, updated_at
	FROM claims
	WHERE id < $1
	ORDER BY id DESC
	LIMIT $2
`

func (cr claimsRepository) GetPage(ctx context.Context, ptr uint64, pageSize uint64) ([]domain.Claim, error) {
	rows, err := cr.Conn(ctx).Query(ctx, queryClaimsGetPage, ptr, pageSize)
	if err != nil {
		return nil, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}
	defer rows.Close()

	return cr.getClaimsFromRows(ctx, rows)
}

const queryClaimsGetFirstUserPage = `
	SELECT 
		id, title, description, category, status, latitude, longitude, created_at, updated_at
	FROM claims
	WHERE created_by = $1
	ORDER BY id DESC
	LIMIT $2
`

func (cr claimsRepository) GetFirstUserPage(ctx context.Context, pageSize uint64, userID uuid.UUID) ([]domain.Claim, error) {
	rows, err := cr.Conn(ctx).Query(ctx, queryClaimsGetFirstUserPage, userID, pageSize)
	if err != nil {
		return nil, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}
	defer rows.Close()

	return cr.getClaimsFromRows(ctx, rows)
}

const queryClaimsGetUserPage = `
	SELECT 
		id, title, description, category, status, latitude, longitude, created_at, updated_at
	FROM claims
	WHERE id < $1 and created_by = $2
	ORDER BY id DESC
	LIMIT $3
`

func (cr claimsRepository) GetUserPage(ctx context.Context, ptr uint64, pageSize uint64, userID uuid.UUID) ([]domain.Claim, error) {
	rows, err := cr.Conn(ctx).Query(ctx, queryClaimsGetUserPage, ptr, userID, pageSize)
	if err != nil {
		return nil, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}
	defer rows.Close()

	return cr.getClaimsFromRows(ctx, rows)
}

const queryClaimsGetByArea = `
	SELECT 
		id, latitude, longitude
	FROM claims
	WHERE 
		latitude <= $1 and longitude >= $2 and latitude > $3 and longitude < $4 and (status = $5 OR (status = $6 and status_updated_at >= $7))
`

// TODO: Пересмотреть работу с поиском активных заявок

func (cr claimsRepository) GetByArea(ctx context.Context, lat1, long1, lat2, long2 float64, acceptedStatus, completedStatus domain.ClaimStatus, startingFrom time.Time) ([]domain.Claim, error) {
	rows, err := cr.Conn(ctx).Query(
		ctx,
		queryClaimsGetByArea,
		lat1,
		long1,
		lat2,
		long2,
		acceptedStatus.String(),
		completedStatus.String(),
		startingFrom,
	)
	if err != nil {
		return nil, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}
	defer rows.Close()

	claims, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.Claim])
	if err != nil {
		return nil, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	return claims, nil
}
