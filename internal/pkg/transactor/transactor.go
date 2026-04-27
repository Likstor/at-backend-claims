package transactor

import (
	"at-backend-claims/internal/pkg/apperror"
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewTransactor(pool *pgxpool.Pool) *Transactor {
	return &Transactor{
		Pool: pool,
	}
}

type Transactor struct {
	Pool *pgxpool.Pool
}

func (t Transactor) WithinTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error {
	tx, err := t.Pool.Begin(ctx)
	if err != nil {
		return apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	slog.InfoContext(ctx, "transaction started")

	ctx = InjectTx(ctx, tx)

	txFuncErr := txFunc(ctx)
	if errors.Is(txFuncErr, apperror.ErrRepository) {
		if err := tx.Rollback(ctx); err != nil {
			return apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
		}

		slog.InfoContext(ctx, "transaction canceled")

		return txFuncErr
	}

	if err := tx.Commit(ctx); err != nil {
		slog.InfoContext(ctx, "transaction canceled")

		return apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	slog.InfoContext(ctx, "transaction committed")

	return txFuncErr
}
