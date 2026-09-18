package handlers_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/models"
)

// TestMalformedJSONBody_Returns422 cobre o branch de erro de decode do JSON
// (corpo que nem chega a ser um objeto válido) em cada rota que lê um body —
// diferente de "faltam campos", que é validado depois, na camada de serviço.
func TestMalformedJSONBody_Returns422(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleGestor}
	id := uuid.New()

	cases := []struct {
		name   string
		method string
		path   string
	}{
		{"solicitations.Create", http.MethodPost, "/api/v1/solicitations"},
		{"solicitations.PatchDraft", http.MethodPatch, "/api/v1/solicitations/" + id.String()},
		{"solicitations.Reject", http.MethodPost, "/api/v1/solicitations/" + id.String() + "/reject"},
		{"solicitations.PatchAnalysis", http.MethodPatch, "/api/v1/solicitations/" + id.String() + "/analysis"},
		{"solicitations.Finalize", http.MethodPost, "/api/v1/solicitations/" + id.String() + "/finalize"},
		{"categories.Create", http.MethodPost, "/api/v1/categories"},
		{"categories.Update", http.MethodPatch, "/api/v1/categories/" + id.String()},
		{"categories.SetApprovers", http.MethodPatch, "/api/v1/categories/" + id.String() + "/approvers"},
		{"users.Create", http.MethodPost, "/api/v1/users"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, d.router.URL+tc.path, bytes.NewReader([]byte("isto não é json")))
			require.NoError(t, err)
			req.Header.Set("Authorization", d.bearerFor(t, actor))

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		})
	}
}

func TestAuthLogin_MalformedJSONBody(t *testing.T) {
	d := newTestServer(t)

	req, err := http.NewRequest(http.MethodPost, d.router.URL+"/api/v1/auth/login", bytes.NewReader([]byte("isto não é json")))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

// TestInvalidIDParam_Returns422 cobre o branch de parseID em rotas ainda não
// exercitadas por outros testes (Get e Submit/Approve já aparecem em outros
// arquivos, mas o restante da família de rotas com {id} não).
func TestInvalidIDParam_Returns422(t *testing.T) {
	d := newTestServer(t)
	actor := models.User{ID: uuid.New(), Role: models.RoleGestor}

	cases := []struct {
		name   string
		method string
		path   string
	}{
		{"solicitations.PatchDraft", http.MethodPatch, "/api/v1/solicitations/not-a-uuid"},
		{"solicitations.Submit", http.MethodPost, "/api/v1/solicitations/not-a-uuid/submit"},
		{"solicitations.Approve", http.MethodPost, "/api/v1/solicitations/not-a-uuid/approve"},
		{"solicitations.Reject", http.MethodPost, "/api/v1/solicitations/not-a-uuid/reject"},
		{"solicitations.PatchAnalysis", http.MethodPatch, "/api/v1/solicitations/not-a-uuid/analysis"},
		{"solicitations.Finalize", http.MethodPost, "/api/v1/solicitations/not-a-uuid/finalize"},
		{"solicitations.GetHistory", http.MethodGet, "/api/v1/solicitations/not-a-uuid/history"},
		{"categories.Update", http.MethodPatch, "/api/v1/categories/not-a-uuid"},
		{"categories.SetApprovers", http.MethodPatch, "/api/v1/categories/not-a-uuid/approvers"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, d.router.URL+tc.path, bytes.NewReader([]byte("{}")))
			require.NoError(t, err)
			req.Header.Set("Authorization", d.bearerFor(t, actor))

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		})
	}
}
