//go:build integration

package migrations_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/migrations"
)

func TestApply_CreatesSchemaAndIsIdempotent(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL não definida — pulando testes de integração")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	// Outros pacotes de teste já aplicaram a migration neste mesmo banco
	// compartilhado, então sem resetar o schema aqui `Apply` só veria
	// `applied=true` e nunca exercitaria de fato o caminho que lê o .sql e
	// cria as tabelas. Zerar o schema garante que o teste exercite o Apply
	// "do zero" de verdade, não só a checagem de idempotência.
	_, err = pool.Exec(ctx, `DROP SCHEMA public CASCADE; CREATE SCHEMA public;`)
	require.NoError(t, err)

	require.NoError(t, migrations.Apply(ctx, pool))

	var tableCount int
	require.NoError(t, pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name = 'solicitations'`).Scan(&tableCount))
	require.Equal(t, 1, tableCount)

	// Rodar de novo não deve falhar nem tentar recriar as tabelas.
	require.NoError(t, migrations.Apply(ctx, pool))

	var migrationCount int
	require.NoError(t, pool.QueryRow(ctx, `SELECT COUNT(*) FROM schema_migrations`).Scan(&migrationCount))
	require.Equal(t, 1, migrationCount)
}
