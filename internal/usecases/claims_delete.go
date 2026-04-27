package usecase

import (
	"at-backend-claims/internal/domain"
	"at-backend-claims/internal/pkg/apperror"
	"at-backend-claims/internal/pkg/reqctx"
	"at-backend-claims/internal/pkg/roles"
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
)

func (cs claimsUsecase) deletePhotos(ctx context.Context, claim domain.Claim) {
	dir := fmt.Sprintf("claims/%d", claim.ID)

	for _, filename := range claim.Photos {
		if err := cs.fileManager.Delete(ctx, filepath.Join(dir, filename)); err != nil {
			slog.ErrorContext(ctx, err.Error())
		}
	}
}

func (cs claimsUsecase) Delete(ctx context.Context, id uint64) error {
	userID, ok := reqctx.GetUserID(ctx)
	if !ok {
		return apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrCtxEmptyUserID)
	}

	userRole, ok := reqctx.GetRole(ctx)
	if !ok {
		return apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrCtxEmptyRole)
	}

	if err := cs.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		claim, err := cs.repo.GetByID(ctx, id)
		if err != nil {
			return err
		}

		// Заявка может быть удалена её владельцем в случае, если она находится на стадии обработки,
		// либо администратором/оператором системы
		if (claim.Status != domain.ClaimStatusPending || claim.CreatedBy != userID) && roles.IsRoleALowerThanB(userRole, roles.Operator) {
			slog.WarnContext(ctx, apperror.ErrUserCannotPerformOperation.Error())
			return apperror.ErrUserCannotPerformOperation
		}

		if err := cs.repo.Delete(ctx, id); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
