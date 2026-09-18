// Package repository define os contratos de acesso a dados usados pela
// camada de serviço. Manter isso como interfaces puras (sem tipos do pgx
// vazando) é o que permite testar SolicitationService com mocks, sem banco.
package repository

import (
	"context"

	"github.com/google/uuid"

	"ekaizen-backend/internal/models"
)

// Transactor demarca uma transação de banco. fn recebe um ctx "marcado" que
// todas as chamadas de repositório dentro dele devem propagar, para que
// participem da mesma transação.
type Transactor interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	ListProfiles(ctx context.Context) ([]models.UserProfile, error)
}

type CategoryRepository interface {
	List(ctx context.Context) ([]models.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	// GetApproverByOrder retorna o usuário aprovador de uma categoria numa
	// determinada ordem (1 ou 2), ou nil se não houver.
	GetApproverByOrder(ctx context.Context, categoryID uuid.UUID, order int16) (*models.User, error)
	// GetApproverOrder retorna a ordem (1 ou 2) de um usuário numa categoria,
	// ou nil se ele não for aprovador dela.
	GetApproverOrder(ctx context.Context, categoryID, userID uuid.UUID) (*int16, error)
}

type SolicitationRepository interface {
	Create(ctx context.Context, s *models.Solicitation) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Solicitation, error)
	// GetForUpdate faz SELECT ... FOR UPDATE — só deve ser chamado dentro de
	// uma transação (Transactor.WithinTx), para travar a linha durante uma
	// transição de estado e evitar corrida em cliques duplicados.
	GetForUpdate(ctx context.Context, id uuid.UUID) (*models.Solicitation, error)
	Update(ctx context.Context, s *models.Solicitation) error
	List(ctx context.Context, filter models.SolicitationFilter) ([]models.Solicitation, error)
	CountsByStatus(ctx context.Context) (map[models.Status]int, error)
	OldestPending(ctx context.Context, limit int) ([]models.OldestPendingItem, error)
	AwaitingApprovalByUser(ctx context.Context) ([]models.UserPendingCount, error)
}

type HistoryRepository interface {
	Create(ctx context.Context, h *models.HistoryEntry) error
	ListBySolicitation(ctx context.Context, solicitationID uuid.UUID) ([]models.HistoryEntry, error)
}
