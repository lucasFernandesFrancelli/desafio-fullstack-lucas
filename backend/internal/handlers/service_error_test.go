package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/models"
)

func TestAuthProfiles_ServiceError(t *testing.T) {
	d := newTestServer(t)
	d.users.On("ListProfiles", mock.Anything).Return(nil, errors.New("falha no banco"))

	resp, err := http.Get(d.router.URL + "/api/v1/auth/profiles")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestAuthLogin_UserNotFound(t *testing.T) {
	d := newTestServer(t)
	missingID := uuid.New()
	d.users.On("GetByID", mock.Anything, missingID).Return(nil, apperrors.NotFound("usuário não encontrado"))

	body, _ := json.Marshal(map[string]string{"userId": missingID.String()})
	resp, err := http.Post(d.router.URL+"/api/v1/auth/login", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAuthMe_ServiceError(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}
	d.users.On("GetByID", mock.Anything, actor.ID).Return(nil, apperrors.NotFound("usuário não encontrado"))

	req, _ := http.NewRequest(http.MethodGet, d.router.URL+"/api/v1/auth/me", nil)
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestUsers_List_ServiceError(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleGestor}
	d.users.On("List", mock.Anything).Return(nil, errors.New("falha no banco"))

	req, _ := http.NewRequest(http.MethodGet, d.router.URL+"/api/v1/users", nil)
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
