//go:build integration

// Roda o seed contra um Postgres real (TEST_DATABASE_URL) e confere que os
// dados fictícios batem com o que o README promete, e que rodar duas vezes
// não duplica nada (idempotência é o que permite o seed rodar em todo boot
// do backend em produção).
package seed_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/seed"
	"ekaizen-backend/migrations"
)

func TestRun_SeedsExpectedDataAndIsIdempotent(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL não definida — pulando testes de integração")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	require.NoError(t, migrations.Apply(ctx, pool))
	_, err = pool.Exec(ctx, `TRUNCATE solicitation_history, solicitations, category_approvers, categories, users CASCADE`)
	require.NoError(t, err)

	require.NoError(t, seed.Run(ctx, pool))

	var userCount, categoryCount, approverCount, solicitationCount, historyCount int
	require.NoError(t, pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&userCount))
	require.NoError(t, pool.QueryRow(ctx, `SELECT COUNT(*) FROM categories`).Scan(&categoryCount))
	require.NoError(t, pool.QueryRow(ctx, `SELECT COUNT(*) FROM category_approvers`).Scan(&approverCount))
	require.NoError(t, pool.QueryRow(ctx, `SELECT COUNT(*) FROM solicitations`).Scan(&solicitationCount))
	require.NoError(t, pool.QueryRow(ctx, `SELECT COUNT(*) FROM solicitation_history`).Scan(&historyCount))

	require.Equal(t, 14, userCount)
	require.Equal(t, 4, categoryCount)
	require.Equal(t, 8, approverCount)
	require.Equal(t, 6, solicitationCount)
	require.Greater(t, historyCount, 6)

	// Rodar de novo não deve duplicar nada (idempotência).
	require.NoError(t, seed.Run(ctx, pool))

	var userCountAgain int
	require.NoError(t, pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&userCountAgain))
	require.Equal(t, userCount, userCountAgain)
}
