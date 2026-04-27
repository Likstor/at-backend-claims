package postgres

import (
	"at-backend-claims/internal/domain"
	"at-backend-claims/internal/pkg/transactor"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dummyClaim = domain.Claim{}

type claimsRepository struct {
	transactor.RepositoryWithTransactor
}

func NewClaimsRepository(pool *pgxpool.Pool) *claimsRepository {
	return &claimsRepository{
		RepositoryWithTransactor: *transactor.NewRepositoryWithTransactor(pool),
	}
}

