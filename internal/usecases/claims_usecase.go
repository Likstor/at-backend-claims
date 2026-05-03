package usecase

import (
	"at-backend-claims/internal/domain"
	"context"
	"io"
	"time"

	"github.com/google/uuid"
)

var dummyClaim = domain.Claim{}

type claimsRepo interface {
	GetByID(ctx context.Context, id uint64) (domain.Claim, error)
	GetFirstPage(ctx context.Context, pageSize uint64) ([]domain.Claim, error)
	GetPage(ctx context.Context, ptr uint64, pageSize uint64) ([]domain.Claim, error)
	GetFirstUserPage(ctx context.Context, pageSize uint64, uid uuid.UUID) ([]domain.Claim, error)
	GetUserPage(ctx context.Context, ptr uint64, pageSize uint64, uid uuid.UUID) ([]domain.Claim, error)
	GetByArea(ctx context.Context, lat1, long1, lat2, long2 float64, createdBy uuid.UUID, pendingStatus, acceptedStatus, completedStatus domain.ClaimStatus, startingFrom time.Time) ([]domain.Claim, error)

	Create(ctx context.Context, data domain.Claim) (uint64, error)

	Delete(ctx context.Context, id uint64) error

	Update(ctx context.Context, data domain.Claim) error
	ChangeStatus(ctx context.Context, id uint64, status domain.ClaimStatus) error
	AddFeedback(ctx context.Context, id uint64, feedback string) error

	WithinTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error
}

type categoryChecker interface {
	SubcategoryExist(ctx context.Context, subcategory string) bool
}

type pointChecker interface {
	IsPointInPolygon(float64, float64) bool
}

type fileManager interface {
	Save(ctx context.Context, filePath string, file io.Reader) error
	Delete(ctx context.Context, filePath string) error
	GetURLToFile(ctx context.Context, filePath string) (string, error)
}

type claimsUsecase struct {
	repo            claimsRepo
	fileManager     fileManager
	categoryChecker categoryChecker
	pointChecker    pointChecker

	maxPageSize                  uint64
	hideCompletedClaimsOlderThan time.Duration
}

func NewClaimsUsecase(
	repo claimsRepo,
	maxPageSize uint64,
	categoryChecker categoryChecker,
	pointChecker pointChecker,
	filesManager fileManager,
	hideCompletedClaimsOlderThan time.Duration,
) *claimsUsecase {
	return &claimsUsecase{
		repo:            repo,
		fileManager:     filesManager,
		maxPageSize:     maxPageSize,
		categoryChecker: categoryChecker,
		pointChecker:    pointChecker,
	}
}
