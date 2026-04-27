package transactor

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryWithTransactor interface {
	Conn(ctx context.Context) Executor
	WithinTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error
}

type repositoryWithTransactor struct {
	Transactor
	pool *pgxpool.Pool
}

func NewRepositoryWithTransactor(pool *pgxpool.Pool) *repositoryWithTransactor {
	return &repositoryWithTransactor{
		Transactor: *NewTransactor(pool),
		pool:       pool,
	}
}

func (repo repositoryWithTransactor) Conn(ctx context.Context) Executor {
	if tx := ExtractTx(ctx); tx != nil {
		return tx
	}

	return repo.pool
}
