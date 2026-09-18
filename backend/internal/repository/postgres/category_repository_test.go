//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/repository/postgres"
)

func TestCategoryRepository_List_WithApprovers(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	repo := postgres.NewCategoryRepository(postgres.NewStore(pool))

	catID := uuid.New()
	approver1, approver2 := uuid.New(), uuid.New()

	_, err := pool.Exec(ctx, `INSERT INTO users (id, name, email, role) VALUES ($1,'Aprovador Um','a1@ekaizen.example','colaborador'), ($2,'Aprovador Dois','a2@ekaizen.example','colaborador')`,
		approver1, approver2)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO categories (id, name, description) VALUES ($1,'Segurança','Descrição')`, catID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO category_approvers (id, category_id, user_id, approval_order) VALUES ($1,$2,$3,1),($4,$2,$5,2)`,
		uuid.New(), catID, approver1, uuid.New(), approver2)
	require.NoError(t, err)

	categories, err := repo.List(ctx)
	require.NoError(t, err)
	require.Len(t, categories, 1)
	require.Len(t, categories[0].Approvers, 2)
	require.Equal(t, "Aprovador Um", categories[0].Approvers[0].UserName)
	require.EqualValues(t, 1, categories[0].Approvers[0].Order)
}

func TestCategoryRepository_GetByID(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	repo := postgres.NewCategoryRepository(postgres.NewStore(pool))

	catID := uuid.New()
	_, err := pool.Exec(ctx, `INSERT INTO categories (id, name) VALUES ($1,'Meio Ambiente')`, catID)
	require.NoError(t, err)

	cat, err := repo.GetByID(ctx, catID)
	require.NoError(t, err)
	require.Equal(t, "Meio Ambiente", cat.Name)

	_, err = repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestCategoryRepository_GetApproverByOrder(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	repo := postgres.NewCategoryRepository(postgres.NewStore(pool))

	catID := uuid.New()
	approverID := uuid.New()
	_, err := pool.Exec(ctx, `INSERT INTO users (id, name, email, role) VALUES ($1,'Aprovador Dois','ap2@ekaizen.example','colaborador')`, approverID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO categories (id, name) VALUES ($1,'Ambiente')`, catID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO category_approvers (id, category_id, user_id, approval_order) VALUES ($1,$2,$3,2)`,
		uuid.New(), catID, approverID)
	require.NoError(t, err)

	approver, err := repo.GetApproverByOrder(ctx, catID, 2)
	require.NoError(t, err)
	require.Equal(t, "Aprovador Dois", approver.Name)

	_, err = repo.GetApproverByOrder(ctx, catID, 1)
	require.Error(t, err)
}

func TestCategoryRepository_GetApproverOrder(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()
	repo := postgres.NewCategoryRepository(postgres.NewStore(pool))

	catID := uuid.New()
	approverID := uuid.New()
	nonApproverID := uuid.New()

	_, err := pool.Exec(ctx, `INSERT INTO users (id, name, email, role) VALUES ($1,'Aprovador','ap@ekaizen.example','colaborador'), ($2,'Outro','outro@ekaizen.example','colaborador')`,
		approverID, nonApproverID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO categories (id, name) VALUES ($1,'Qualidade')`, catID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO category_approvers (id, category_id, user_id, approval_order) VALUES ($1,$2,$3,2)`,
		uuid.New(), catID, approverID)
	require.NoError(t, err)

	order, err := repo.GetApproverOrder(ctx, catID, approverID)
	require.NoError(t, err)
	require.NotNil(t, order)
	require.EqualValues(t, 2, *order)

	order, err = repo.GetApproverOrder(ctx, catID, nonApproverID)
	require.NoError(t, err)
	require.Nil(t, order)
}

func TestCategoryRepository_DuplicateApprovalOrder_ViolatesConstraint(t *testing.T) {
	pool := setupPool(t)
	ctx := context.Background()

	catID := uuid.New()
	u1, u2 := uuid.New(), uuid.New()
	_, err := pool.Exec(ctx, `INSERT INTO users (id, name, email, role) VALUES ($1,'U1','u1@ekaizen.example','colaborador'), ($2,'U2','u2@ekaizen.example','colaborador')`, u1, u2)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO categories (id, name) VALUES ($1,'Produtividade')`, catID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO category_approvers (id, category_id, user_id, approval_order) VALUES ($1,$2,$3,1)`, uuid.New(), catID, u1)
	require.NoError(t, err)

	// mesma categoria + mesma ordem (1) para um segundo usuário deve violar
	// a constraint UNIQUE(category_id, approval_order) — a regra de negócio
	// "dois aprovadores distintos, ordem previamente definida" é garantida
	// pelo próprio schema, não só pela aplicação.
	_, err = pool.Exec(ctx, `INSERT INTO category_approvers (id, category_id, user_id, approval_order) VALUES ($1,$2,$3,1)`, uuid.New(), catID, u2)
	require.Error(t, err)
}
