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

func TestSolicitations_List_Success(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}

	d.sols.On("List", mock.Anything, mock.Anything).Return([]models.Solicitation{
		{ID: uuid.New(), Title: "Item", RequesterID: actor.ID, Status: models.StatusEmAnalise},
	}, nil)

	req, _ := http.NewRequest(http.MethodGet, d.router.URL+"/api/v1/solicitations?status=em_analise", nil)
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var items []models.SolicitationSummary
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&items))
	require.Len(t, items, 1)
}

func TestSolicitations_List_InvalidCategoryID(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}

	req, _ := http.NewRequest(http.MethodGet, d.router.URL+"/api/v1/solicitations?categoryId=not-a-uuid", nil)
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestSolicitations_PatchDraft_Success(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, RequesterID: actor.ID, Status: models.StatusRascunho, Title: "Antigo"}

	d.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)
	d.sols.On("Update", mock.Anything, mock.Anything).Return(nil)
	d.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)

	body, _ := json.Marshal(map[string]string{"title": "Novo título"})
	req, _ := http.NewRequest(http.MethodPatch, d.router.URL+"/api/v1/solicitations/"+solID.String(), bytes.NewReader(body))
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSolicitations_PatchAnalysis_Success(t *testing.T) {
	d := newTestServer(t)
	analyst := models.User{ID: uuid.New(), Role: models.RoleAnalista}
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, Status: models.StatusEmAnalise}

	d.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)
	d.sols.On("Update", mock.Anything, mock.Anything).Return(nil)
	d.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)

	sev := 4
	body, _ := json.Marshal(map[string]int{"severity": sev})
	req, _ := http.NewRequest(http.MethodPatch, d.router.URL+"/api/v1/solicitations/"+solID.String()+"/analysis", bytes.NewReader(body))
	req.Header.Set("Authorization", d.bearerFor(t, analyst))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSolicitations_Finalize_Success(t *testing.T) {
	d := newTestServer(t)
	analyst := models.User{ID: uuid.New(), Role: models.RoleAnalista}
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, Status: models.StatusEmAnalise}

	d.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)
	d.sols.On("Update", mock.Anything, mock.Anything).Return(nil)
	d.hist.On("Create", mock.Anything, mock.Anything).Return(nil)
	d.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)

	body, _ := json.Marshal(map[string]any{
		"severity": 4, "urgency": 3, "trend": 5, "analysisNotes": "Parecer final",
	})
	req, _ := http.NewRequest(http.MethodPost, d.router.URL+"/api/v1/solicitations/"+solID.String()+"/finalize", bytes.NewReader(body))
	req.Header.Set("Authorization", d.bearerFor(t, analyst))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSolicitations_Finalize_Forbidden(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}
	solID := uuid.New()

	req, _ := http.NewRequest(http.MethodPost, d.router.URL+"/api/v1/solicitations/"+solID.String()+"/finalize", bytes.NewReader([]byte("{}")))
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestSolicitations_GetHistory_Success(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}
	solID := uuid.New()

	d.sols.On("GetByID", mock.Anything, solID).Return(&models.Solicitation{ID: solID, RequesterID: actor.ID, Status: models.StatusEmAnalise}, nil)
	d.hist.On("ListBySolicitation", mock.Anything, solID).Return([]models.HistoryEntry{}, nil)

	req, _ := http.NewRequest(http.MethodGet, d.router.URL+"/api/v1/solicitations/"+solID.String()+"/history", nil)
	req.Header.Set("Authorization", d.bearerFor(t, actor))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSolicitations_GetHistory_HiddenDraft(t *testing.T) {
	d := newTestServer(t)
	other := models.User{ID: uuid.New(), Role: models.RoleColaborador}
	solID := uuid.New()

	d.sols.On("GetByID", mock.Anything, solID).Return(&models.Solicitation{ID: solID, RequesterID: uuid.New(), Status: models.StatusRascunho}, nil)

	req, _ := http.NewRequest(http.MethodGet, d.router.URL+"/api/v1/solicitations/"+solID.String()+"/history", nil)
	req.Header.Set("Authorization", d.bearerFor(t, other))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
