// Package migrations aplica os arquivos *.up.sql deste diretório em ordem,
// registrando cada versão numa tabela schema_migrations. É deliberadamente
// simples (sem downs em runtime) porque roda automaticamente no boot do
// backend — a avaliação é via link público, ninguém vai rodar uma
// ferramenta de migration manualmente. Os arquivos .down.sql continuam no
// repositório para rollback manual documentado no README.
package migrations

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed *.up.sql
var files embed.FS

func Apply(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`); err != nil {
		return fmt.Errorf("criar tabela schema_migrations: %w", err)
	}

	entries, err := fs.Glob(files, "*.up.sql")
	if err != nil {
		return err
	}
	sort.Strings(entries)

	for _, name := range entries {
		version := strings.TrimSuffix(name, ".up.sql")

		var applied bool
		if err := pool.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`, version).Scan(&applied); err != nil {
			return err
		}
		if applied {
			continue
		}

		contents, err := files.ReadFile(name)
		if err != nil {
			return err
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, string(contents)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("aplicar migration %s: %w", name, err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, version); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
		fmt.Printf("migrations: aplicada %s\n", name)
	}
	return nil
}
