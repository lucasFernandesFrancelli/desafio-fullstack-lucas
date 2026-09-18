package services

import (
	"context"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/models"
	"ekaizen-backend/internal/repository"
)

// DashboardService responde à necessidade do gestor de enxergar onde o
// processo está parado sem perguntar individualmente para cada pessoa.
type DashboardService struct {
	solicitations repository.SolicitationRepository
}

func NewDashboardService(solicitations repository.SolicitationRepository) *DashboardService {
	return &DashboardService{solicitations: solicitations}
}

const oldestPendingLimit = 10

func (s *DashboardService) GetSummary(ctx context.Context, actor models.User) (*models.DashboardSummary, error) {
	if actor.Role != models.RoleGestor {
		return nil, apperrors.Forbidden("apenas o gestor pode acessar o dashboard")
	}

	counts, err := s.solicitations.CountsByStatus(ctx)
	if err != nil {
		return nil, err
	}
	oldest, err := s.solicitations.OldestPending(ctx, oldestPendingLimit)
	if err != nil {
		return nil, err
	}
	awaiting, err := s.solicitations.AwaitingApprovalByUser(ctx)
	if err != nil {
		return nil, err
	}
	if oldest == nil {
		oldest = []models.OldestPendingItem{}
	}
	if awaiting == nil {
		awaiting = []models.UserPendingCount{}
	}

	return &models.DashboardSummary{
		CountsByStatus: counts,
		OldestPending:  oldest,
		AwaitingByUser: awaiting,
	}, nil
}
