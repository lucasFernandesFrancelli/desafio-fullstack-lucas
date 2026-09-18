package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"ekaizen-backend/internal/seed"
	"ekaizen-backend/migrations"
)

// runSeed conecta, aplica migrations e roda o seed — extraído de main() para
// ser testável sem chamar os.Exit indiretamente via log.Fatalf.
func runSeed(ctx context.Context, databaseURL string) error {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("conectar ao banco: %w", err)
	}
	defer pool.Close()

	if err := migrations.Apply(ctx, pool); err != nil {
		return fmt.Errorf("aplicar migrations: %w", err)
	}
	if err := seed.Run(ctx, pool); err != nil {
		return fmt.Errorf("aplicar seed: %w", err)
	}
	return nil
}
