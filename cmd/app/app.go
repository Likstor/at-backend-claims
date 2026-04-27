package main

import (
	"at-backend-claims/internal/handlers"
	"at-backend-claims/internal/handlers/middleware"
	"at-backend-claims/internal/infrastructure/db/postgres"
	filemanager "at-backend-claims/internal/infrastructure/external/file_manager"
	"at-backend-claims/internal/pkg/config"
	"at-backend-claims/internal/pkg/districts"
	"at-backend-claims/internal/pkg/health"
	"at-backend-claims/internal/pkg/logs"
	"at-backend-claims/internal/pkg/pgclient"
	"at-backend-claims/internal/pkg/tokens"
	usecase "at-backend-claims/internal/usecases"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	handler := logs.NewHandlerMiddleware(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))
	slog.SetDefault(slog.New(handler))

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	cfg := Config{}
	if err := config.Load(&cfg); err != nil {
		slog.Error(err.Error())
		return
	}

	pool, err := pgclient.NewClient(ctx1, 5, pgclient.Config{
		Host:     cfg.StorageHost,
		Port:     cfg.StoragePort,
		Database: cfg.StorageName,
		Username: cfg.StorageUser,
		Password: cfg.StoragePassword,
		SSLMode:  cfg.StorageSSLMode,
	})
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer pool.Close()

	tokensVerifier := tokens.NewTokensVerifier(
		15*time.Minute,
		cfg.JWKSURI,
		cfg.TokensIssuer,
		cfg.TokensKeyID,
		5*time.Second,
	)

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	filesManager, err := filemanager.NewS3Client(
		ctx2,
		cfg.ObjectStorageAccessKeyID,
		cfg.ObjectStorageSecretAccessKey,
		"",
		cfg.ObjectStorageBucket,
		cfg.ObjectStorageEndpoint,
		cfg.ObjectStorageAccessHost,
	)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	tverCentral, err := districts.LoadPolygonFromFile("./resources/polygons/tver.central")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	districtsService := districts.NewDistrictsService(map[string]districts.Polygon{
		"tver.central": tverCentral,
	})

	claimsRepo := postgres.NewClaimsRepository(pool)

	categoriesRepo := postgres.NewCategoriesRepository(pool)

	ctx3, cancel3 := context.WithCancel(context.Background())
	defer cancel3()

	categoriesUsecase := usecase.NewCategoriesUsecase(
		ctx3,
		categoriesRepo,
		100,
	)

	claimsUsecase := usecase.NewClaimsUsecase(
		claimsRepo,
		100,
		categoriesUsecase,
		districtsService,
		filesManager,
		time.Hour*24*14,
	)

	claimsHandler := handlers.NewClaimsHandler(claimsUsecase)
	claimsHandlerForAdmins := handlers.NewClaimsHandlerForAdmins(claimsUsecase)

	categoriesHandler := handlers.NewCategoriesHandler(categoriesUsecase)
	categoriesHandlerForAdmins := handlers.NewCategoriesHandlerForAdmins(categoriesUsecase)

	districtsHandler := handlers.NewDistrictsHandler(districtsService)

	mainMux := http.NewServeMux()

	claimsHandler.Setup("/claims", tokensVerifier.VerifyJWT, mainMux)
	claimsHandlerForAdmins.Setup("/admins/claims", tokensVerifier.VerifyJWT, mainMux)
	categoriesHandler.Setup("/categories", tokensVerifier.VerifyJWT, mainMux)
	categoriesHandlerForAdmins.Setup("/admins/categories", tokensVerifier.VerifyJWT, mainMux)
	districtsHandler.Setup("/districts", tokensVerifier.VerifyJWT, mainMux)

	mainMuxWrapped := middleware.Cors(mainMux, cfg.AllowOriginCors)
	mainMuxWrapped = middleware.Logger(mainMux)
	mainMuxWrapped = middleware.Correlation(mainMuxWrapped)

	finMux := http.NewServeMux()
	finMux.Handle("/v1/", mainMuxWrapped)
	health.Setup(finMux)

	server := &http.Server{
		Addr:              ":" + cfg.BackendPort,
		Handler:           finMux,
		ReadTimeout:       cfg.BackendReadTimeout,
		ReadHeaderTimeout: cfg.BackendReadHeaderTimeout,
		WriteTimeout:      cfg.BackendWriteTimeout,
		IdleTimeout:       cfg.BackendIdleTimeout,
	}

	run(context.Background(), server)

	slog.Info("app closed")
}

func run(ctx context.Context, server *http.Server) {
	var wg sync.WaitGroup
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	wg.Add(1)
	go func() {
		slog.Info("HTTP server started")
		err := server.ListenAndServe()
		if err == http.ErrServerClosed {
			slog.Info("http server closed")
			wg.Done()
			return
		}

		slog.Error("http server closed with error: " + err.Error())
	}()

	wg.Add(1)
	go func() {
		<-ctx.Done()

		err := server.Shutdown(ctx)
		if err != nil {
			slog.Error("http server shutdown error: " + err.Error())
		}

		wg.Done()
	}()

	wg.Wait()
}
