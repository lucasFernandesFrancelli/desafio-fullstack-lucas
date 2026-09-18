package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"ekaizen-backend/internal/auth"
	"ekaizen-backend/internal/config"
	"ekaizen-backend/internal/handlers"
	"ekaizen-backend/internal/repository/postgres"
	"ekaizen-backend/internal/seed"
	"ekaizen-backend/internal/server"
	"ekaizen-backend/internal/services"
	"ekaizen-backend/migrations"
)

// buildApp conecta no banco, aplica migrations e seed, e monta o router com
// todas as dependências. Extraído de main() para que um teste de integração
// possa exercitar a mesma montagem sem precisar chamar ListenAndServe.
func buildApp(ctx context.Context, cfg *config.Config) (http.Handler, *pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("conectar ao banco: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("ping no banco: %w", err)
	}
	if err := migrations.Apply(ctx, pool); err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("aplicar migrations: %w", err)
	}
	if err := seed.Run(ctx, pool); err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("aplicar seed: %w", err)
	}

	store := postgres.NewStore(pool)
	userRepo := postgres.NewUserRepository(store)
	categoryRepo := postgres.NewCategoryRepository(store)
	solicitationRepo := postgres.NewSolicitationRepository(store)
	historyRepo := postgres.NewHistoryRepository(store)

	issuer := auth.NewJWTIssuer(cfg.JWTSecret)

	authService := services.NewAuthService(userRepo, issuer)
	solicitationService := services.NewSolicitationService(store, solicitationRepo, categoryRepo, historyRepo)
	dashboardService := services.NewDashboardService(solicitationRepo)
	adminService := services.NewAdminService(store, userRepo, categoryRepo)

	router := server.NewRouter(server.Dependencies{
		Issuer:              issuer,
		CORSOrigin:          cfg.CORSOrigin,
		AuthHandler:         handlers.NewAuthHandler(authService),
		CategoryHandler:     handlers.NewCategoryHandler(categoryRepo, adminService),
		SolicitationHandler: handlers.NewSolicitationHandler(solicitationService),
		DashboardHandler:    handlers.NewDashboardHandler(dashboardService),
		UserHandler:         handlers.NewUserHandler(adminService),
	})

	return router, pool, nil
}
