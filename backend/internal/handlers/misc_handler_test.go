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

	"ekaizen-backend/internal/models"
)

func TestHealthz(t *testing.T) {
	d := newTestServer(t)

	resp, err := http.Get(d.router.URL + "/healthz")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthProfiles_PublicNoToken(t *testing.T) {
	d := newTestServer(t)
	d.users.On("ListProfiles", mock.Anything).Return([]models.UserProfile{}, nil)

	resp, err := http.Get(d.router.URL + "/api/v1/auth/profiles")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthLogin_Success(t *testing.T) {
	d := newTestServer(t)
	user := &models.User{ID: uuid.New(), Name: "Ana", Role: models.RoleColaborador}
	d.users.On("GetByID", mock.Anything, user.ID).Return(user, nil)

	body, _ := json.Marshal(map[string]string{"userId": user.ID.String()})
	resp, err := http.Post(d.router.URL+"/api/v1/auth/login", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var payload struct {
		Token string `json:"token"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&payload))
	require.NotEmpty(t, payload.Token)
}

func TestAuthLogin_InvalidUUID(t *testing.T) {
	d := newTestServer(t)

	body, _ := json.Marshal(map[string]string{"userId": "não-é-um-uuid"})
	resp, err := http.Post(d.router.URL+"/api/v1/auth/login", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestAuthMe_Unauthorized(t *testing.T) {
	d := newTestServer(t)

	resp, err := http.Get(d.router.URL + "/api/v1/auth/me")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthMe_Success(t *testing.T) {
	d := newTestServer(t)
	user := models.User{ID: uuid.New(), Name: "Ana", Role: models.RoleColaborador}
	d.users.On("GetByID", mock.Anything, user.ID).Return(&user, nil)

	req, _ := http.NewRequest(http.MethodGet, d.router.URL+"/api/v1/auth/me", nil)
	req.Header.Set("Authorization", d.bearerFor(t, user))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCategories_RequireAuth(t *testing.T) {
	d := newTestServer(t)

	resp, err := http.Get(d.router.URL + "/api/v1/categories")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthMe_InvalidToken(t *testing.T) {
	d := newTestServer(t)

	req, _ := http.NewRequest(http.MethodGet, d.router.URL+"/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer isto-nao-e-um-token-valido")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestCategories_ListRepositoryError(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}
	d.cats.On("List", mock.Anything).Return(nil, errors.New("db down"))

	req, _ := http.NewRequest(http.MethodGet, d.router.URL+"/api/v1/categories", nil)
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestCategories_ListWithAuth(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}
	d.cats.On("List", mock.Anything).Return([]models.Category{{ID: uuid.New(), Name: "Qualidade"}}, nil)

	req, _ := http.NewRequest(http.MethodGet, d.router.URL+"/api/v1/categories", nil)
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
