package usecase

import (
	"at-backend-claims/internal/domain"
	"at-backend-claims/internal/pkg/apperror"
	"context"
)

func (cs claimsUsecase) Update(ctx context.Context, updatedClaim domain.Claim) error {
	updatedClaim.Status = domain.ClaimStatusPending

	if !cs.pointChecker.IsPointInPolygon(updatedClaim.Longitude, updatedClaim.Latitude) {
		return apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrPointIsNotInPolygon)
	}

	if err := cs.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		claim, err := cs.repo.GetByID(ctx, updatedClaim.ID)
		if err != nil {
			return err
		}

		if updatedClaim.CreatedBy != claim.CreatedBy || claim.Status != domain.ClaimStatusPending {
			return apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrUserCannotPerformOperation)
		}

		if err := cs.repo.Update(ctx, updatedClaim); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (cs claimsUsecase) AddFeedback(ctx context.Context, id uint64, feedback string) error {
	if err := cs.repo.AddFeedback(ctx, id, feedback); err != nil {
		return err
	}

	return nil
}

const opClaimsChangeStatus = "usecase.Claims.ChangeStatus"

func (cs claimsUsecase) ChangeStatus(ctx context.Context, id uint64, status domain.ClaimStatus) error {
	if status == domain.ClaimStatusUnknown {
		return apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrUnknownClaimStatus)
	}

	if err := cs.repo.ChangeStatus(ctx, id, status); err != nil {
		return err
	}

	return nil
}
