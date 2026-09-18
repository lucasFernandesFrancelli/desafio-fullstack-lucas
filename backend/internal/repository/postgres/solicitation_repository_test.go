//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/models"
	"ekaizen-backend/internal/repository/postgres"
)

func seedRequester(t *testing.T, ctx context.Context, store *postgres.Store) uuid.UUID {
	t.Helper()
	id := uuid.New()
	_, err := store.Pool.Exec(ctx, `INSERT INTO users (id, name, email, role) VALUES ($1,'Solicitante','sol@ekaizen.example','colaborador')`, id)
	require.NoError(t, err)
	return id
}

func TestSolicitationRepository_CreateAndGetByID(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	store := postgres.NewStore(pool)
	repo := postgres.NewSolicitationRepository(store)
	requesterID := seedRequester(t, ctx, store)

	sol := &models.Solicitation{
		ID: uuid.New(), Title: "Piso escorregadio", RequesterID: requesterID, Status: models.StatusRascunho,
	}
	require.NoError(t, repo.Create(ctx, sol))
	require.False(t, sol.CreatedAt.IsZero())

	fetched, err := repo.GetByID(ctx, sol.ID)
	require.NoError(t, err)
	require.Equal(t, "Piso escorregadio", fetched.Title)
	require.Equal(t, "Solicitante", fetched.RequesterName)
	require.Equal(t, models.StatusRascunho, fetched.Status)
}

func TestSolicitationRepository_GetByID_NotFound(t *testing.T) {
	pool := setupPool(t)
	repo := postgres.NewSolicitationRepository(postgres.NewStore(pool))

	_, err := repo.GetByID(context.Background(), uuid.New())

	appErr, ok := apperrors.As(err)
	require.True(t, ok)
	require.Equal(t, apperrors.CodeNotFound, appErr.Code)
}

func TestSolicitationRepository_GetForUpdate_NotFound(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	store := postgres.NewStore(pool)
	repo := postgres.NewSolicitationRepository(store)

	err := store.WithinTx(ctx, func(txCtx context.Context) error {
		_, err := repo.GetForUpdate(txCtx, uuid.New())
		return err
	})

	appErr, ok := apperrors.As(err)
	require.True(t, ok)
	require.Equal(t, apperrors.CodeNotFound, appErr.Code)
}

func TestSolicitationRepository_GetByID_AttachesCategoryName(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	store := postgres.NewStore(pool)
	repo := postgres.NewSolicitationRepository(store)
	requesterID := seedRequester(t, ctx, store)

	catID := uuid.New()
	_, err := pool.Exec(ctx, `INSERT INTO categories (id, name) VALUES ($1,'Categoria com nome')`, catID)
	require.NoError(t, err)

	sol := &models.Solicitation{ID: uuid.New(), Title: "Item", RequesterID: requesterID, CategoryID: &catID, Status: models.StatusRascunho}
	require.NoError(t, repo.Create(ctx, sol))

	fetched, err := repo.GetByID(ctx, sol.ID)
	require.NoError(t, err)
	require.Equal(t, "Categoria com nome", fetched.CategoryName)
}

func TestSolicitationRepository_GetByID_PropagatesGenericError(t *testing.T) {
	pool := setupPool(t)
	repo := postgres.NewSolicitationRepository(postgres.NewStore(pool))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := repo.GetByID(ctx, uuid.New())

	require.Error(t, err)
	_, ok := apperrors.As(err)
	require.False(t, ok, "erro genérico de contexto cancelado não deve virar apperrors.NotFound")
}

func TestSolicitationRepository_GetForUpdate_PropagatesGenericError(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	store := postgres.NewStore(pool)
	repo := postgres.NewSolicitationRepository(store)

	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()

	err := store.WithinTx(ctx, func(txCtx context.Context) error {
		_, err := repo.GetForUpdate(cancelCtx, uuid.New())
		return err
	})

	require.Error(t, err)
}

func TestSolicitationRepository_Create_PropagatesError(t *testing.T) {
	pool := setupPool(t)
	repo := postgres.NewSolicitationRepository(postgres.NewStore(pool))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repo.Create(ctx, &models.Solicitation{ID: uuid.New(), Title: "X", RequesterID: uuid.New(), Status: models.StatusRascunho})

	require.Error(t, err)
}

func TestSolicitationRepository_List_PropagatesQueryError(t *testing.T) {
	pool := setupPool(t)
	repo := postgres.NewSolicitationRepository(postgres.NewStore(pool))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := repo.List(ctx, models.SolicitationFilter{})

	require.Error(t, err)
}

func TestSolicitationRepository_CountsByStatus_PropagatesQueryError(t *testing.T) {
	pool := setupPool(t)
	repo := postgres.NewSolicitationRepository(postgres.NewStore(pool))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := repo.CountsByStatus(ctx)

	require.Error(t, err)
}

func TestSolicitationRepository_OldestPending_PropagatesQueryError(t *testing.T) {
	pool := setupPool(t)
	repo := postgres.NewSolicitationRepository(postgres.NewStore(pool))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := repo.OldestPending(ctx, 10)

	require.Error(t, err)
}

func TestSolicitationRepository_AwaitingApprovalByUser_PropagatesQueryError(t *testing.T) {
	pool := setupPool(t)
	repo := postgres.NewSolicitationRepository(postgres.NewStore(pool))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := repo.AwaitingApprovalByUser(ctx)

	require.Error(t, err)
}

func TestSolicitationRepository_Update_ComputesGeneratedPriority(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	store := postgres.NewStore(pool)
	repo := postgres.NewSolicitationRepository(store)
	requesterID := seedRequester(t, ctx, store)

	sol := &models.Solicitation{ID: uuid.New(), Title: "Item", RequesterID: requesterID, Status: models.StatusEmAnalise}
	require.NoError(t, repo.Create(ctx, sol))

	sev, urg, trend := int16(4), int16(3), int16(5)
	sol.Severity, sol.Urgency, sol.Trend = &sev, &urg, &trend
	sol.AnalysisNotes = "Parecer"
	sol.Status = models.StatusFinalizada

	require.NoError(t, repo.Update(ctx, sol))
	require.NotNil(t, sol.Priority)
	require.EqualValues(t, 60, *sol.Priority) // 4 * 3 * 5, calculado pela coluna GENERATED do Postgres

	fetched, err := repo.GetByID(ctx, sol.ID)
	require.NoError(t, err)
	require.EqualValues(t, 60, *fetched.Priority)
}

func TestSolicitationRepository_List_FilterByStatusAndQuery(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	store := postgres.NewStore(pool)
	repo := postgres.NewSolicitationRepository(store)
	requesterID := seedRequester(t, ctx, store)

	sol1 := &models.Solicitation{ID: uuid.New(), Title: "Vazamento de óleo na prensa", RequesterID: requesterID, Status: models.StatusRascunho}
	sol2 := &models.Solicitation{ID: uuid.New(), Title: "Troca de EPI obrigatório", RequesterID: requesterID, Status: models.StatusEmAnalise}
	require.NoError(t, repo.Create(ctx, sol1))
	require.NoError(t, repo.Create(ctx, sol2))

	rascunho := models.StatusRascunho
	results, err := repo.List(ctx, models.SolicitationFilter{Status: &rascunho})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Equal(t, sol1.ID, results[0].ID)

	results, err = repo.List(ctx, models.SolicitationFilter{Query: "vazamento"})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Equal(t, sol1.ID, results[0].ID)
}

func TestSolicitationRepository_CountsByStatus(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	store := postgres.NewStore(pool)
	repo := postgres.NewSolicitationRepository(store)
	requesterID := seedRequester(t, ctx, store)

	require.NoError(t, repo.Create(ctx, &models.Solicitation{ID: uuid.New(), Title: "A", RequesterID: requesterID, Status: models.StatusRascunho}))
	require.NoError(t, repo.Create(ctx, &models.Solicitation{ID: uuid.New(), Title: "B", RequesterID: requesterID, Status: models.StatusRascunho}))

	counts, err := repo.CountsByStatus(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, counts[models.StatusRascunho])
	require.Equal(t, 0, counts[models.StatusFinalizada])
}

func TestSolicitationRepository_OldestPending(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	store := postgres.NewStore(pool)
	repo := postgres.NewSolicitationRepository(store)
	requesterID := seedRequester(t, ctx, store)

	require.NoError(t, repo.Create(ctx, &models.Solicitation{ID: uuid.New(), Title: "Parado", RequesterID: requesterID, Status: models.StatusRascunho}))
	require.NoError(t, repo.Create(ctx, &models.Solicitation{ID: uuid.New(), Title: "Finalizado", RequesterID: requesterID, Status: models.StatusFinalizada}))

	items, err := repo.OldestPending(ctx, 10)
	require.NoError(t, err)
	require.Len(t, items, 1) // finalizada não conta como pendente
	require.Equal(t, "Parado", items[0].Title)
}

func TestSolicitationRepository_AwaitingApprovalByUser(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	store := postgres.NewStore(pool)
	repo := postgres.NewSolicitationRepository(store)
	requesterID := seedRequester(t, ctx, store)

	approverID := uuid.New()
	catID := uuid.New()
	_, err := pool.Exec(ctx, `INSERT INTO users (id, name, email, role) VALUES ($1,'Aprovador','aprov@ekaizen.example','colaborador')`, approverID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO categories (id, name) VALUES ($1,'Categoria X')`, catID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO category_approvers (id, category_id, user_id, approval_order) VALUES ($1,$2,$3,1)`,
		uuid.New(), catID, approverID)
	require.NoError(t, err)

	step1 := int16(1)
	sol := &models.Solicitation{ID: uuid.New(), Title: "Aguardando", RequesterID: requesterID, Status: models.StatusEmAprovacao, CategoryID: &catID}
	require.NoError(t, repo.Create(ctx, sol))
	sol.CurrentApprovalStep = &step1
	require.NoError(t, repo.Update(ctx, sol))

	counts, err := repo.AwaitingApprovalByUser(ctx)
	require.NoError(t, err)
	require.Len(t, counts, 1)
	require.Equal(t, "Aprovador", counts[0].UserName)
	require.Equal(t, 1, counts[0].Count)
}

func TestSolicitationRepository_GetForUpdate_WithinTransaction(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	store := postgres.NewStore(pool)
	repo := postgres.NewSolicitationRepository(store)
	requesterID := seedRequester(t, ctx, store)

	sol := &models.Solicitation{ID: uuid.New(), Title: "Item travado", RequesterID: requesterID, Status: models.StatusRascunho}
	require.NoError(t, repo.Create(ctx, sol))

	err := store.WithinTx(ctx, func(txCtx context.Context) error {
		locked, err := repo.GetForUpdate(txCtx, sol.ID)
		require.NoError(t, err)
		require.Equal(t, sol.ID, locked.ID)
		return nil
	})
	require.NoError(t, err)
}
