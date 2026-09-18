//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/models"
	"ekaizen-backend/internal/repository/postgres"
)

func TestHistoryRepository_CreateAndListInOrder(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	store := postgres.NewStore(pool)
	solRepo := postgres.NewSolicitationRepository(store)
	histRepo := postgres.NewHistoryRepository(store)
	requesterID := seedRequester(t, ctx, store)

	sol := &models.Solicitation{ID: uuid.New(), Title: "Item", RequesterID: requesterID, Status: models.StatusRascunho}
	require.NoError(t, solRepo.Create(ctx, sol))

	rascunho := models.StatusRascunho
	emAprovacao := models.StatusEmAprovacao

	require.NoError(t, histRepo.Create(ctx, &models.HistoryEntry{
		ID: uuid.New(), SolicitationID: sol.ID, ToStatus: models.StatusRascunho,
		Action: models.ActionCriada, ActorID: requesterID,
	}))
	require.NoError(t, histRepo.Create(ctx, &models.HistoryEntry{
		ID: uuid.New(), SolicitationID: sol.ID, FromStatus: &rascunho, ToStatus: emAprovacao,
		Action: models.ActionEnviada, ActorID: requesterID,
	}))

	entries, err := histRepo.ListBySolicitation(ctx, sol.ID)
	require.NoError(t, err)
	require.Len(t, entries, 2)
	require.Equal(t, models.ActionCriada, entries[0].Action)
	require.Equal(t, models.ActionEnviada, entries[1].Action)
	require.Equal(t, "Solicitante", entries[0].ActorName)
}

func TestHistoryRepository_ListBySolicitation_PropagatesQueryError(t *testing.T) {
	pool := setupPool(t)
	repo := postgres.NewHistoryRepository(postgres.NewStore(pool))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := repo.ListBySolicitation(ctx, uuid.New())

	require.Error(t, err)
}

func TestHistoryRepository_Create_PropagatesError(t *testing.T) {
	pool := setupPool(t)
	repo := postgres.NewHistoryRepository(postgres.NewStore(pool))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repo.Create(ctx, &models.HistoryEntry{
		ID: uuid.New(), SolicitationID: uuid.New(), ToStatus: models.StatusRascunho,
		Action: models.ActionCriada, ActorID: uuid.New(),
	})

	require.Error(t, err)
}
