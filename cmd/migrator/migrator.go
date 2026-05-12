package main

import (
	"at-backend-claims/internal/pkg/config"
	"at-backend-claims/internal/pkg/logs"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	handler := logs.NewHandlerMiddleware(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	slog.SetDefault(slog.New(handler))

	cfg := Config{}
	if err := config.Load(&cfg); err != nil {
		slog.Error(err.Error())
		return
	}

	storageURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?x-migrations-table=%s&sslmode=%s",
		cfg.StorageUser,
		cfg.StoragePassword,
		cfg.StorageHost,
		cfg.StoragePort,
		cfg.StorageName,
		cfg.MigrationsTable,
		cfg.StorageSSLMode,
	)

	migration, err := migrate.New("file://"+cfg.MigrationsPath, storageURL)
	if err != nil {
		slog.Error("migration create error")

		return
	}
	defer migration.Close()

	ok := false

	for i := 1; i <= 5; i++ {
		slog.Info(fmt.Sprintf("attempt %v to migration", i))
		if err := migration.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				slog.Info("no new migrations")
				ok = true
				return
			}

			slog.Error("migrations apply error")
			continue
		}

		slog.Info("migrations applied")
		return
	}

	if !ok {
		slog.Error("migrations not applied")
		os.Exit(1)
	}
}
