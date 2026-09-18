package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/models"
	"ekaizen-backend/internal/services"
	"ekaizen-backend/internal/testutil"
)

func TestDashboardService_Forbidden_NotManager(t *testing.T) {
	sols := &testutil.SolicitationRepo{}
	svc := services.NewDashboardService(sols)

	_, err := svc.GetSummary(context.Background(), models.User{ID: uuid.New(), Role: models.RoleColaborador})

	appErr, ok := apperrors.As(err)
	require.True(t, ok)
	require.Equal(t, apperrors.CodeForbidden, appErr.Code)
}

func TestDashboardService_Success(t *testing.T) {
	sols := &testutil.SolicitationRepo{}
	svc := services.NewDashboardService(sols)

	counts := map[models.Status]int{models.StatusRascunho: 2}
	oldest := []models.OldestPendingItem{{ID: uuid.New(), Title: "Item parado"}}
	awaiting := []models.UserPendingCount{{UserID: uuid.New(), UserName: "Fulano", Count: 3}}

	sols.On("CountsByStatus", mock.Anything).Return(counts, nil)
	sols.On("OldestPending", mock.Anything, 10).Return(oldest, nil)
	sols.On("AwaitingApprovalByUser", mock.Anything).Return(awaiting, nil)

	summary, err := svc.GetSummary(context.Background(), models.User{ID: uuid.New(), Role: models.RoleGestor})

	require.NoError(t, err)
	require.Equal(t, 2, summary.CountsByStatus[models.StatusRascunho])
	require.Len(t, summary.OldestPending, 1)
	require.Len(t, summary.AwaitingByUser, 1)
}

func TestDashboardService_PropagatesRepositoryError(t *testing.T) {
	sols := &testutil.SolicitationRepo{}
	svc := services.NewDashboardService(sols)

	sols.On("CountsByStatus", mock.Anything).Return(nil, errors.New("falha no banco"))

	_, err := svc.GetSummary(context.Background(), models.User{ID: uuid.New(), Role: models.RoleGestor})

	require.Error(t, err)
}
