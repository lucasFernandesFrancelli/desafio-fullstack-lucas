package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ctxKey string

const txCtxKey ctxKey = "pgx-tx"

// Store guarda o pool de conexões e implementa repository.Transactor. Cada
// repositório concreto recebe um *Store e usa Store.db(ctx) para obter a
// conexão certa (a transação ativa, se houver, ou o pool).
type Store struct {
	Pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{Pool: pool}
}

func (s *Store) db(ctx context.Context) DBTX {
	if tx, ok := ctx.Value(txCtxKey).(pgx.Tx); ok {
		return tx
	}
	return s.Pool
}

// WithinTx abre uma transação, injeta no ctx passado para fn, e comita ao
// final. Qualquer erro retornado por fn provoca rollback.
func (s *Store) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }() // no-op após commit bem-sucedido

	txCtx := context.WithValue(ctx, txCtxKey, tx)
	if err := fn(txCtx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
