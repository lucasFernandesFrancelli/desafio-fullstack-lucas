//go:build integration

// Testes de integração de repositório: rodam contra um Postgres real
// apontado por TEST_DATABASE_URL (sem testcontainers — este ambiente não
// tem Docker disponível, então usamos um Postgres local/CI de verdade, o
// que também funciona perfeitamente em GitHub Actions via um serviço
// postgres). Rode com: go test ./... -tags=integration
package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/migrations"
)

func setupPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL não definida — pulando testes de integração")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	require.NoError(t, migrations.Apply(context.Background(), pool))

	_, err = pool.Exec(context.Background(),
		`TRUNCATE solicitation_history, solicitations, category_approvers, categories, users CASCADE`)
	require.NoError(t, err)

	return pool
}
