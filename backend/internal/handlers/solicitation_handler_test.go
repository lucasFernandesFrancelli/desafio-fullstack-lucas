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

func TestSolicitations_Create_Success(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}

	d.sols.On("Create", mock.Anything, mock.Anything).Return(nil)
	d.hist.On("Create", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(map[string]string{"title": "Piso escorregadio no refeitório"})
	req, _ := http.NewRequest(http.MethodPost, d.router.URL+"/api/v1/solicitations", bytes.NewReader(body))
	req.Header.Set("Authorization", d.bearerFor(t, actor))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestSolicitations_Create_MissingTitle(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}

	body, _ := json.Marshal(map[string]string{})
	req, _ := http.NewRequest(http.MethodPost, d.router.URL+"/api/v1/solicitations", bytes.NewReader(body))
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestSolicitations_Get_HiddenDraftReturnsNotFound(t *testing.T) {
	d := newTestServer(t)
	solID := uuid.New()
	other := models.User{ID: uuid.New(), Role: models.RoleColaborador}

	d.sols.On("GetByID", mock.Anything, solID).
		Return(&models.Solicitation{ID: solID, RequesterID: uuid.New(), Status: models.StatusRascunho}, nil)

	req, _ := http.NewRequest(http.MethodGet, d.router.URL+"/api/v1/solicitations/"+solID.String(), nil)
	req.Header.Set("Authorization", d.bearerFor(t, other))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestSolicitations_Submit_ValidationError(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}
	solID := uuid.New()

	d.sols.On("GetForUpdate", mock.Anything, solID).
		Return(&models.Solicitation{ID: solID, RequesterID: actor.ID, Status: models.StatusRascunho, Title: "Só título"}, nil)

	req, _ := http.NewRequest(http.MethodPost, d.router.URL+"/api/v1/solicitations/"+solID.String()+"/submit", nil)
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	var payload struct {
		Error struct {
			Fields map[string]string `json:"fields"`
		} `json:"error"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&payload))
	require.Contains(t, payload.Error.Fields, "categoryId")
}

func TestSolicitations_Approve_ForbiddenWrongTurn(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}
	catID := uuid.New()
	solID := uuid.New()
	step1 := int16(1)

	d.sols.On("GetForUpdate", mock.Anything, solID).
		Return(&models.Solicitation{ID: solID, Status: models.StatusEmAprovacao, CategoryID: &catID, CurrentApprovalStep: &step1}, nil)
	d.cats.On("GetApproverOrder", mock.Anything, catID, actor.ID).Return(nil, nil)

	req, _ := http.NewRequest(http.MethodPost, d.router.URL+"/api/v1/solicitations/"+solID.String()+"/approve", nil)
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestSolicitations_Reject_MissingReason(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}
	solID := uuid.New()

	body, _ := json.Marshal(map[string]string{"reason": ""})
	req, _ := http.NewRequest(http.MethodPost, d.router.URL+"/api/v1/solicitations/"+solID.String()+"/reject", bytes.NewReader(body))
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestSolicitations_Get_InvalidID(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}

	req, _ := http.NewRequest(http.MethodGet, d.router.URL+"/api/v1/solicitations/not-a-uuid", nil)
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}
