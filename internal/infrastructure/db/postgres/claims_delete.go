package postgres

import (
	"at-backend-claims/internal/pkg/apperror"
	"context"
)

const queryClaimsDelete = `
	DELETE FROM claims
	WHERE id = $1
`

func (cr claimsRepository) Delete(ctx context.Context, id uint64) error {
	resp, err := cr.Conn(ctx).Exec(ctx, queryClaimsDelete, id)
	if err != nil {
		return apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	if resp.RowsAffected() < 1 {
		return apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrClaimNotExists)
	}

	return nil
}
