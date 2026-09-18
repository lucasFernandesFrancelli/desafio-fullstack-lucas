// Package testutil fornece mocks (testify/mock) das interfaces de
// repositório, compartilhados pelos testes unitários de internal/services e
// pelos testes de handler (httptest) de internal/handlers.
package testutil

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"ekaizen-backend/internal/models"
)

func asPtr[T any](v any) *T {
	if v == nil {
		return nil
	}
	return v.(*T)
}

func asSlice[T any](v any) []T {
	if v == nil {
		return nil
	}
	return v.([]T)
}

func asMap[K comparable, V any](v any) map[K]V {
	if v == nil {
		return nil
	}
	return v.(map[K]V)
}

// Transactor executa fn direto, sem transação real — suficiente para testar
// a lógica de negócio isolada do banco.
type Transactor struct{}

func (t *Transactor) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type UserRepo struct{ mock.Mock }

func (m *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	return asPtr[models.User](args.Get(0)), args.Error(1)
}

func (m *UserRepo) ListProfiles(ctx context.Context) ([]models.UserProfile, error) {
	args := m.Called(ctx)
	return asSlice[models.UserProfile](args.Get(0)), args.Error(1)
}

func (m *UserRepo) List(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return asSlice[models.User](args.Get(0)), args.Error(1)
}

func (m *UserRepo) Create(ctx context.Context, u *models.User) error {
	return m.Called(ctx, u).Error(0)
}

type CategoryRepo struct{ mock.Mock }

func (m *CategoryRepo) List(ctx context.Context) ([]models.Category, error) {
	args := m.Called(ctx)
	return asSlice[models.Category](args.Get(0)), args.Error(1)
}

func (m *CategoryRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	args := m.Called(ctx, id)
	return asPtr[models.Category](args.Get(0)), args.Error(1)
}

func (m *CategoryRepo) GetApproverByOrder(ctx context.Context, categoryID uuid.UUID, order int16) (*models.User, error) {
	args := m.Called(ctx, categoryID, order)
	return asPtr[models.User](args.Get(0)), args.Error(1)
}

func (m *CategoryRepo) GetApproverOrder(ctx context.Context, categoryID, userID uuid.UUID) (*int16, error) {
	args := m.Called(ctx, categoryID, userID)
	return asPtr[int16](args.Get(0)), args.Error(1)
}

func (m *CategoryRepo) Create(ctx context.Context, c *models.Category) error {
	return m.Called(ctx, c).Error(0)
}

func (m *CategoryRepo) UpdateInfo(ctx context.Context, id uuid.UUID, name, description string) error {
	return m.Called(ctx, id, name, description).Error(0)
}

func (m *CategoryRepo) SetApprovers(ctx context.Context, categoryID, firstApproverID, secondApproverID uuid.UUID) error {
	return m.Called(ctx, categoryID, firstApproverID, secondApproverID).Error(0)
}

type SolicitationRepo struct{ mock.Mock }

func (m *SolicitationRepo) Create(ctx context.Context, s *models.Solicitation) error {
	return m.Called(ctx, s).Error(0)
}

func (m *SolicitationRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Solicitation, error) {
	args := m.Called(ctx, id)
	return asPtr[models.Solicitation](args.Get(0)), args.Error(1)
}

func (m *SolicitationRepo) GetForUpdate(ctx context.Context, id uuid.UUID) (*models.Solicitation, error) {
	args := m.Called(ctx, id)
	return asPtr[models.Solicitation](args.Get(0)), args.Error(1)
}

func (m *SolicitationRepo) Update(ctx context.Context, s *models.Solicitation) error {
	return m.Called(ctx, s).Error(0)
}

func (m *SolicitationRepo) List(ctx context.Context, filter models.SolicitationFilter) ([]models.Solicitation, error) {
	args := m.Called(ctx, filter)
	return asSlice[models.Solicitation](args.Get(0)), args.Error(1)
}

func (m *SolicitationRepo) CountsByStatus(ctx context.Context) (map[models.Status]int, error) {
	args := m.Called(ctx)
	return asMap[models.Status, int](args.Get(0)), args.Error(1)
}

func (m *SolicitationRepo) OldestPending(ctx context.Context, limit int) ([]models.OldestPendingItem, error) {
	args := m.Called(ctx, limit)
	return asSlice[models.OldestPendingItem](args.Get(0)), args.Error(1)
}

func (m *SolicitationRepo) AwaitingApprovalByUser(ctx context.Context) ([]models.UserPendingCount, error) {
	args := m.Called(ctx)
	return asSlice[models.UserPendingCount](args.Get(0)), args.Error(1)
}

type HistoryRepo struct{ mock.Mock }

func (m *HistoryRepo) Create(ctx context.Context, h *models.HistoryEntry) error {
	return m.Called(ctx, h).Error(0)
}

func (m *HistoryRepo) ListBySolicitation(ctx context.Context, solicitationID uuid.UUID) ([]models.HistoryEntry, error) {
	args := m.Called(ctx, solicitationID)
	return asSlice[models.HistoryEntry](args.Get(0)), args.Error(1)
}
