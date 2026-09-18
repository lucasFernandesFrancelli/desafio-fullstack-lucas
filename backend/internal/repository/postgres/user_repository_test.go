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

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	pool := setupPool(t)
	repo := postgres.NewUserRepository(postgres.NewStore(pool))

	_, err := repo.GetByID(context.Background(), uuid.New())

	appErr, ok := apperrors.As(err)
	require.True(t, ok)
	require.Equal(t, apperrors.CodeNotFound, appErr.Code)
}

func TestUserRepository_GetByID_Found(t *testing.T) {
	pool := setupPool(t)
	repo := postgres.NewUserRepository(postgres.NewStore(pool))
	ctx := context.Background()

	id := uuid.New()
	_, err := pool.Exec(ctx, `INSERT INTO users (id, name, email, role) VALUES ($1,$2,$3,$4)`,
		id, "Ana Teste", "ana.teste@ekaizen.example", models.RoleColaborador)
	require.NoError(t, err)

	user, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, "Ana Teste", user.Name)
	require.Equal(t, models.RoleColaborador, user.Role)
}

func TestUserRepository_ListProfiles_IncludesApproverFor(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	repo := postgres.NewUserRepository(postgres.NewStore(pool))

	userID := uuid.New()
	catID := uuid.New()
	_, err := pool.Exec(ctx, `INSERT INTO users (id, name, email, role) VALUES ($1,$2,$3,$4)`,
		userID, "Ricardo Aprovador", "ricardo.aprovador@ekaizen.example", models.RoleColaborador)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO categories (id, name) VALUES ($1,$2)`, catID, "Categoria Teste")
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO category_approvers (id, category_id, user_id, approval_order) VALUES ($1,$2,$3,1)`,
		uuid.New(), catID, userID)
	require.NoError(t, err)

	profiles, err := repo.ListProfiles(ctx)
	require.NoError(t, err)
	require.Len(t, profiles, 1)
	require.Len(t, profiles[0].ApproverFor, 1)
	require.Equal(t, "Categoria Teste", profiles[0].ApproverFor[0].CategoryName)
	require.EqualValues(t, 1, profiles[0].ApproverFor[0].Order)
}

func TestUserRepository_Create_Success(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	repo := postgres.NewUserRepository(postgres.NewStore(pool))

	user := &models.User{ID: uuid.New(), Name: "Pessoa Nova", Email: "pessoa.nova@ekaizen.example", Role: models.RoleAnalista}
	require.NoError(t, repo.Create(ctx, user))
	require.False(t, user.CreatedAt.IsZero())

	fetched, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, "Pessoa Nova", fetched.Name)
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	repo := postgres.NewUserRepository(postgres.NewStore(pool))

	email := "duplicado@ekaizen.example"
	require.NoError(t, repo.Create(ctx, &models.User{ID: uuid.New(), Name: "Um", Email: email, Role: models.RoleColaborador}))

	err := repo.Create(ctx, &models.User{ID: uuid.New(), Name: "Dois", Email: email, Role: models.RoleColaborador})
	appErr, ok := apperrors.As(err)
	require.True(t, ok)
	require.Equal(t, apperrors.CodeValidation, appErr.Code)
}

func TestUserRepository_List_ReturnsAllOrderedByName(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	repo := postgres.NewUserRepository(postgres.NewStore(pool))

	require.NoError(t, repo.Create(ctx, &models.User{ID: uuid.New(), Name: "Zeca", Email: "zeca@ekaizen.example", Role: models.RoleColaborador}))
	require.NoError(t, repo.Create(ctx, &models.User{ID: uuid.New(), Name: "Amanda", Email: "amanda@ekaizen.example", Role: models.RoleColaborador}))

	users, err := repo.List(ctx)
	require.NoError(t, err)
	require.Len(t, users, 2)
	require.Equal(t, "Amanda", users[0].Name)
	require.Equal(t, "Zeca", users[1].Name)
}
