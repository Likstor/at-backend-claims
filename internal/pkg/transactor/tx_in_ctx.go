package transactor

import (
	"context"

	"github.com/jackc/pgx/v5"
)

const txCtxKey = "executor"

func InjectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txCtxKey, tx)
}

func ExtractTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txCtxKey).(pgx.Tx); ok {
		return tx
	}

	return nil
}