package handlers_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/models"
)

func TestDashboard_ForbiddenForColaborador(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}

	req, _ := http.NewRequest(http.MethodGet, d.router.URL+"/api/v1/dashboard", nil)
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestDashboard_SuccessForGestor(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleGestor}

	d.sols.On("CountsByStatus", mock.Anything).Return(map[models.Status]int{models.StatusRascunho: 1}, nil)
	d.sols.On("OldestPending", mock.Anything, mock.Anything).Return([]models.OldestPendingItem{}, nil)
	d.sols.On("AwaitingApprovalByUser", mock.Anything).Return([]models.UserPendingCount{}, nil)

	req, _ := http.NewRequest(http.MethodGet, d.router.URL+"/api/v1/dashboard", nil)
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
