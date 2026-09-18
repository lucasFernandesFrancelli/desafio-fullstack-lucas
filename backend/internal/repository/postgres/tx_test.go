//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/models"
	"ekaizen-backend/internal/repository/postgres"
)

func TestStore_WithinTx_RollsBackOnError(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	store := postgres.NewStore(pool)
	repo := postgres.NewSolicitationRepository(store)
	requesterID := seedRequester(t, ctx, store)

	solID := uuid.New()
	expectedErr := errors.New("falha proposital")

	err := store.WithinTx(ctx, func(txCtx context.Context) error {
		sol := &models.Solicitation{ID: solID, Title: "Nunca deve persistir", RequesterID: requesterID, Status: models.StatusRascunho}
		if err := repo.Create(txCtx, sol); err != nil {
			return err
		}
		return expectedErr
	})

	require.ErrorIs(t, err, expectedErr)

	_, getErr := repo.GetByID(ctx, solID)
	require.Error(t, getErr, "a criação deveria ter sido desfeita pelo rollback")
}

func TestStore_WithinTx_PropagatesBeginError(t *testing.T) {
	pool := setupPool(t)
	store := postgres.NewStore(pool)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := store.WithinTx(ctx, func(txCtx context.Context) error {
		t.Fatal("fn não deveria ser chamada quando Begin falha")
		return nil
	})

	require.Error(t, err)
}

func TestStore_WithinTx_CommitsOnSuccess(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	store := postgres.NewStore(pool)
	repo := postgres.NewSolicitationRepository(store)
	requesterID := seedRequester(t, ctx, store)

	solID := uuid.New()

	err := store.WithinTx(ctx, func(txCtx context.Context) error {
		sol := &models.Solicitation{ID: solID, Title: "Deve persistir", RequesterID: requesterID, Status: models.StatusRascunho}
		return repo.Create(txCtx, sol)
	})
	require.NoError(t, err)

	fetched, err := repo.GetByID(ctx, solID)
	require.NoError(t, err)
	require.Equal(t, "Deve persistir", fetched.Title)
}
