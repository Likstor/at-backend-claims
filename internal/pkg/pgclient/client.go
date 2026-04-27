package pgclient

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrFailedConnectionToDatabase = fmt.Errorf("failed connection to database")
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	SSLMode  string
}

func NewClient(ctx context.Context, connAttempts int, config Config) (*pgxpool.Pool, error) {
	slog.InfoContext(ctx, "pool creation started")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.Username, config.Password, config.Host, config.Port, config.Database, config.SSLMode)

	for attempt := 1; attempt <= connAttempts; attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

		pool, err := pgxpool.New(attemptCtx, connString)
		if err != nil {
			cancel()
			slog.WarnContext(ctx, "failed to create pool", "attempt", attempt, "error", err)
			if attempt < connAttempts {
				time.Sleep(5 * time.Second)
			}
			continue
		}

		err = pool.Ping(attemptCtx)
		cancel()
		if err != nil {
			pool.Close()
			slog.WarnContext(ctx, "failed to ping database", "attempt", attempt, "error", err)
			if attempt < connAttempts {
				time.Sleep(5 * time.Second)
			}
			continue
		}

		slog.InfoContext(ctx, "successful connection to the database")
		slog.InfoContext(ctx, "pool creation successful")
		return pool, nil
	}

	slog.ErrorContext(ctx, "pool creation failed after all attempts", "attempts", connAttempts)
	return nil, ErrFailedConnectionToDatabase
}
