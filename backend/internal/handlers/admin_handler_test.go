package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/models"
)

func TestUsers_List_ForbiddenForNonGestor(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}

	req, _ := http.NewRequest(http.MethodGet, d.router.URL+"/api/v1/users", nil)
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestUsers_List_SuccessForGestor(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleGestor}
	d.users.On("List", mock.Anything).Return([]models.User{{ID: uuid.New(), Name: "Fulano"}}, nil)

	req, _ := http.NewRequest(http.MethodGet, d.router.URL+"/api/v1/users", nil)
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUsers_Create_Success(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleGestor}
	d.users.On("Create", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(map[string]string{"name": "Nova Pessoa", "email": "nova@ekaizen.example", "role": "analista"})
	req, _ := http.NewRequest(http.MethodPost, d.router.URL+"/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestUsers_Create_ValidationError(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleGestor}

	body, _ := json.Marshal(map[string]string{})
	req, _ := http.NewRequest(http.MethodPost, d.router.URL+"/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestCategories_Create_ForbiddenForNonGestor(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}

	body, _ := json.Marshal(map[string]string{
		"name": "Ergonomia", "firstApproverId": uuid.New().String(), "secondApproverId": uuid.New().String(),
	})
	req, _ := http.NewRequest(http.MethodPost, d.router.URL+"/api/v1/categories", bytes.NewReader(body))
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestCategories_Create_Success(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleGestor}
	firstID, secondID := uuid.New(), uuid.New()

	d.users.On("GetByID", mock.Anything, firstID).Return(&models.User{ID: firstID}, nil)
	d.users.On("GetByID", mock.Anything, secondID).Return(&models.User{ID: secondID}, nil)
	d.cats.On("Create", mock.Anything, mock.Anything).Return(nil)
	d.cats.On("SetApprovers", mock.Anything, mock.Anything, firstID, secondID).Return(nil)
	d.cats.On("GetByID", mock.Anything, mock.Anything).Return(&models.Category{ID: uuid.New(), Name: "Ergonomia"}, nil)

	body, _ := json.Marshal(map[string]string{
		"name": "Ergonomia", "firstApproverId": firstID.String(), "secondApproverId": secondID.String(),
	})
	req, _ := http.NewRequest(http.MethodPost, d.router.URL+"/api/v1/categories", bytes.NewReader(body))
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestCategories_Update_Success(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleGestor}
	id := uuid.New()

	d.cats.On("GetByID", mock.Anything, id).Return(&models.Category{ID: id, Name: "Antigo"}, nil)
	d.cats.On("UpdateInfo", mock.Anything, id, "Novo Nome", "").Return(nil)

	body, _ := json.Marshal(map[string]string{"name": "Novo Nome"})
	req, _ := http.NewRequest(http.MethodPatch, d.router.URL+"/api/v1/categories/"+id.String(), bytes.NewReader(body))
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCategories_SetApprovers_Success(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleGestor}
	id := uuid.New()
	firstID, secondID := uuid.New(), uuid.New()

	d.users.On("GetByID", mock.Anything, firstID).Return(&models.User{ID: firstID}, nil)
	d.users.On("GetByID", mock.Anything, secondID).Return(&models.User{ID: secondID}, nil)
	d.cats.On("GetByID", mock.Anything, id).Return(&models.Category{ID: id}, nil)
	d.cats.On("SetApprovers", mock.Anything, id, firstID, secondID).Return(nil)

	body, _ := json.Marshal(map[string]string{"firstApproverId": firstID.String(), "secondApproverId": secondID.String()})
	req, _ := http.NewRequest(http.MethodPatch, d.router.URL+"/api/v1/categories/"+id.String()+"/approvers", bytes.NewReader(body))
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCategories_SetApprovers_SameApproverValidation(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleGestor}
	id := uuid.New()
	sameID := uuid.New()

	body, _ := json.Marshal(map[string]string{"firstApproverId": sameID.String(), "secondApproverId": sameID.String()})
	req, _ := http.NewRequest(http.MethodPatch, d.router.URL+"/api/v1/categories/"+id.String()+"/approvers", bytes.NewReader(body))
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}
