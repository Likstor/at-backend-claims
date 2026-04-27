package postgres

import (
	"at-backend-claims/internal/domain"
	"at-backend-claims/internal/pkg/apperror"
	"context"
)

// TODO: Вынести фото в отдельную таблицу

const queryClaimsUpdate = `
	UPDATE claims
	SET 
		title = $1,
		description = $2,
		longitude = $3,
		latitude = $4,
		updated_at = NOW() AT TIME ZONE 'UTC'
	WHERE
		id = $5
`

func (cr claimsRepository) Update(ctx context.Context, data domain.Claim) error {
	resp, err := cr.Conn(ctx).Exec(
		ctx,
		queryClaimsUpdate,
		data.Title,
		data.Description,
		data.Longitude,
		data.Latitude,
		data.ID,
	)
	if err != nil {
		return apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	if resp.RowsAffected() < 1 {
		return apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrClaimNotExists)
	}

	return nil
}

const queryClaimsChangeStatus = `
	UPDATE claims
	SET
		status = $1,
		status_updated_at = NOW() AT TIME ZONE 'UTC'
	WHERE
		id = $2
`

func (cr claimsRepository) ChangeStatus(ctx context.Context, id uint64, status domain.ClaimStatus) error {
	resp, err := cr.Conn(ctx).Exec(ctx, queryClaimsChangeStatus, status.String(), id)
	if err != nil {
		return apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	if resp.RowsAffected() < 1 {
		return apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrClaimNotExists)
	}

	return nil
}

const queryClaimsAddFeedback = `
	UPDATE claims
	SET
		feedback = $1,
		feedback_updated_at = NOW() AT TIME ZONE 'UTC'
	WHERE
		id = $2
`

func (cr claimsRepository) AddFeedback(ctx context.Context, id uint64, feedback string) error {
	resp, err := cr.Conn(ctx).Exec(ctx, queryClaimsAddFeedback, feedback, id)
	if err != nil {
		return apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	if resp.RowsAffected() < 1 {
		return apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrClaimNotExists)
	}
	return nil
}
