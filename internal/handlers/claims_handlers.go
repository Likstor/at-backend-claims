package handlers

import (
	"at-backend-claims/internal/domain"
	"at-backend-claims/internal/handlers/middleware"
	"at-backend-claims/internal/pkg/roles"
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type claimsUsecase interface {
	GetByID(ctx context.Context, id uint64) (domain.Claim, error)
	GetFirstUserPage(ctx context.Context, uid uuid.UUID, pageSize uint64) ([]domain.Claim, error)
	GetUserPage(ctx context.Context, cursor uint64, uid uuid.UUID, pageSize uint64) ([]domain.Claim, error)
	GetByArea(ctx context.Context, lat1, long1, lat2, long2 float64) ([]domain.Claim, error)

	CreateProblem(ctx context.Context, claim domain.Claim, files [][]byte) (uint64, error)
	CreateProposal(ctx context.Context, claim domain.Claim) (uint64, error)

	Delete(ctx context.Context, id uint64) error

	Update(ctx context.Context, updatedClaim domain.Claim) error
}

type claimsHandler struct {
	usecase claimsUsecase

	claimTitleCharactersMaxCount       int
	claimDescriptionCharactersMaxCount int
}

func NewClaimsHandler(uc claimsUsecase) *claimsHandler {
	return &claimsHandler{
		usecase: uc,
	}
}

func (c claimsHandler) Setup(prefix string, verifier middleware.Verifier, mux *http.ServeMux) {
	muxWithAuth := http.NewServeMux()

	muxWithAuth.HandleFunc("POST /problem", c.createProblem)
	muxWithAuth.HandleFunc("POST /proposal", c.createProposal)
	muxWithAuth.HandleFunc("GET /{id}", c.getById)
	muxWithAuth.HandleFunc("PUT /{id}", c.update)
	muxWithAuth.HandleFunc("DELETE /{id}", c.delete)
	
	muxWithAuth.HandleFunc("GET /page", c.getPage)
	muxWithAuth.HandleFunc("GET /map", c.getClaimsByArea)

	muxWithAuthWrapped := middleware.Authorization(muxWithAuth, roles.User, verifier)

	mux.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(prefix, muxWithAuthWrapped))
}

type claimsUsecaseForAdmins interface {
	claimsUsecase

	GetFirstPage(ctx context.Context, pageSize uint64) ([]domain.Claim, error)
	GetPage(ctx context.Context, cursor, pageSize uint64) ([]domain.Claim, error)

	AddFeedback(ctx context.Context, id uint64, feedback string) error
	ChangeStatus(ctx context.Context, id uint64, status domain.ClaimStatus) error
}

type claimsHandlerForAdmins struct {
	usecase claimsUsecaseForAdmins
}

func NewClaimsHandlerForAdmins(uc claimsUsecaseForAdmins) *claimsHandlerForAdmins {
	return &claimsHandlerForAdmins{
		usecase: uc,
	}
}

func (c claimsHandlerForAdmins) Setup(prefix string, verifier middleware.Verifier, mux *http.ServeMux) {
	muxWithAuth := http.NewServeMux()

	muxWithAuth.HandleFunc("GET /{id}", c.getClaimByID)
	muxWithAuth.HandleFunc("GET /page", c.getPageClaims)
	muxWithAuth.HandleFunc("POST /{id}/feedback", c.addFeedbackToClaim)
	muxWithAuth.HandleFunc("POST /{id}/status", c.changeClaimStatus)
	muxWithAuth.HandleFunc("DELETE /{id}", c.deleteClaim)

	muxWithAuth.HandleFunc("GET /users/page", c.getUserClaimsPage)

	muxWithAuthWrapped := middleware.Authorization(muxWithAuth, roles.Operator, verifier)

	mux.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(prefix, muxWithAuthWrapped))
}
