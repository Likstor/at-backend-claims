package usecase

import (
	"at-backend-claims/internal/domain"
	fm "at-backend-claims/internal/infrastructure/external/file_manager"
	"at-backend-claims/internal/pkg/apperror"
	"at-backend-claims/internal/pkg/image"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
)

func (cs claimsUsecase) CreateProblem(ctx context.Context, claim domain.Claim, files [][]byte) (uint64, error) {
	claim.Status = domain.ClaimStatusPending

	if !cs.categoryChecker.SubcategoryExist(ctx, claim.Category) {
		return 0, apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrUnknownClaimCategory)
	}

	if !cs.pointChecker.IsPointInPolygon(claim.Longitude, claim.Latitude) {
		return 0, apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrPointIsNotInPolygon)
	}

	if len(files) != 0 {
		claim.Photos = make([]string, 0, len(files))
	}

	fileCount := 0
	preparedFiles := make(map[string]io.Reader)
	for _, file := range files {
		ok, imageType := image.IsImage(file)
		if !ok {
			slog.WarnContext(ctx, apperror.ErrIsNotImage.Error())
			continue // Пропускаем файл, если он не является jpeg/png
		}

		if imageType == "image/png" {
			bytes, err := image.ConvertToJpeg(file)
			if err != nil {
				slog.WarnContext(ctx, err.Error())
				continue // Пропускаем файл (png), если не получилось конвертировать
			}

			file = bytes
		}

		r := bytes.NewReader(file)

		name := fm.NewFileName(".jpeg")
		preparedFiles[name] = r

		claim.Photos = append(claim.Photos, name)

		if fileCount >= 2 {
			break
		}
		fileCount++
	}

	id, err := cs.repo.Create(ctx, claim)
	if err != nil {
		cs.deletePhotos(ctx, claim)

		return id, err
	}

	for name, file := range preparedFiles {
		dir := fmt.Sprintf("claims/%d", id)

		if err := cs.fileManager.Save(ctx, filepath.Join(dir, name), file); err != nil {
			cs.deletePhotos(ctx, claim)

			if err := cs.repo.Delete(ctx, id); err != nil {
				slog.ErrorContext(ctx, err.Error())
			}

			return 0, err
		}
	}

	return id, nil
}

func (cs claimsUsecase) CreateProposal(ctx context.Context, claim domain.Claim) (uint64, error) {
	claim.Status = domain.ClaimStatusPending
	claim.Category = domain.ProposalCategory

	if !cs.pointChecker.IsPointInPolygon(claim.Longitude, claim.Latitude) {
		return 0, apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrPointIsNotInPolygon)
	}

	id, err := cs.repo.Create(ctx, claim)
	if err != nil {
		return id, err
	}

	return id, nil
}
