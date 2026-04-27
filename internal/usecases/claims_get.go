package usecase

import (
	"at-backend-claims/internal/domain"
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func (cs claimsUsecase) claimPhotosPathToURLs(ctx context.Context, claim *domain.Claim) {
	urls := make([]string, 0, len(claim.Photos))

	dir := fmt.Sprintf("claims/%v", claim.ID)

	for j := 0; j < len(claim.Photos); j++ {
		filePath := filepath.Join(dir, claim.Photos[j])

		photoURL, err := cs.fileManager.GetURLToFile(ctx, filePath)
		if err != nil {
			slog.WarnContext(ctx, err.Error())
			continue
		}

		urls = append(urls, photoURL)
	}

	claim.Photos = urls
}

func (cs claimsUsecase) GetByID(ctx context.Context, id uint64) (domain.Claim, error) {
	claim, err := cs.repo.GetByID(ctx, id)
	if err != nil {
		return dummyClaim, err
	}

	cs.claimPhotosPathToURLs(ctx, &claim)

	return claim, nil
}

func (cs claimsUsecase) GetFirstPage(ctx context.Context, pageSize uint64) ([]domain.Claim, error) {
	if pageSize > cs.maxPageSize {
		pageSize = cs.maxPageSize
	}

	claims, err := cs.repo.GetFirstPage(ctx, pageSize)
	if err != nil {
		return nil, err
	}

	return claims, err
}

func (cs claimsUsecase) GetPage(ctx context.Context, cursor, pageSize uint64) ([]domain.Claim, error) {
	if pageSize > cs.maxPageSize {
		pageSize = cs.maxPageSize
	}

	claims, err := cs.repo.GetPage(ctx, cursor, pageSize)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (cs claimsUsecase) GetFirstUserPage(ctx context.Context, uid uuid.UUID, pageSize uint64) ([]domain.Claim, error) {
	if pageSize > cs.maxPageSize {
		pageSize = cs.maxPageSize
	}

	claims, err := cs.repo.GetFirstUserPage(ctx, pageSize, uid)
	if err != nil {
		return nil, err
	}

	return claims, err
}

func (cs claimsUsecase) GetUserPage(ctx context.Context, cursor uint64, uid uuid.UUID, pageSize uint64) ([]domain.Claim, error) {
	if pageSize > cs.maxPageSize {
		pageSize = cs.maxPageSize
	}

	claims, err := cs.repo.GetUserPage(ctx, cursor, pageSize, uid)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (cs claimsUsecase) GetByArea(ctx context.Context, lat1, long1, lat2, long2 float64) ([]domain.Claim, error) {
	noOlderThan := time.Now().Add(-cs.hideCompletedClaimsOlderThan)

	claims, err := cs.repo.GetByArea(
		ctx,
		lat1,
		long1,
		lat2,
		long2,
		domain.ClaimStatusAccepted,
		domain.ClaimStatusCompleted,
		noOlderThan,
	)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
