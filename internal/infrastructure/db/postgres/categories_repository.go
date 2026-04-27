package postgres

import (
	"at-backend-claims/internal/pkg/transactor"

	"github.com/jackc/pgx/v5/pgxpool"
)

type categoriesRepository struct {
	transactor.RepositoryWithTransactor
}

func NewCategoriesRepository(pool *pgxpool.Pool) *categoriesRepository {
	return &categoriesRepository{
		RepositoryWithTransactor: *transactor.NewRepositoryWithTransactor(pool),
	}
}
