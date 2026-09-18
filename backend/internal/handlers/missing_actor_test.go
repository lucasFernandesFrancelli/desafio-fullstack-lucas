package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/auth"
	"ekaizen-backend/internal/handlers"
	"ekaizen-backend/internal/services"
	"ekaizen-backend/internal/testutil"
)

// TestHandlers_MissingActor_ReturnUnauthorized chama cada handler protegido
// diretamente (sem passar pelo middleware.Auth), simulando o cenário
// defensivo em que, por engano, uma rota fica sem o middleware — o handler
// deve recusar com 401 em vez de continuar com um actor zerado.
func TestHandlers_MissingActor_ReturnUnauthorized(t *testing.T) {
	issuer := auth.NewJWTIssuer("test-secret")
	sols := &testutil.SolicitationRepo{}
	cats := &testutil.CategoryRepo{}
	hist := &testutil.HistoryRepo{}
	users := &testutil.UserRepo{}

	solSvc := services.NewSolicitationService(&testutil.Transactor{}, sols, cats, hist)
	adminSvc := services.NewAdminService(&testutil.Transactor{}, users, cats)
	authSvc := services.NewAuthService(users, issuer)

	solHandler := handlers.NewSolicitationHandler(solSvc)
	catHandler := handlers.NewCategoryHandler(cats, adminSvc)
	userHandler := handlers.NewUserHandler(adminSvc)
	authHandler := handlers.NewAuthHandler(authSvc)
	dashHandler := handlers.NewDashboardHandler(services.NewDashboardService(sols))

	cases := []struct {
		name   string
		fn     func(http.ResponseWriter, *http.Request)
		method string
		body   string
	}{
		{"Solicitation.Create", solHandler.Create, http.MethodPost, "{}"},
		{"Solicitation.List", solHandler.List, http.MethodGet, ""},
		{"Solicitation.Get", solHandler.Get, http.MethodGet, ""},
		{"Solicitation.PatchDraft", solHandler.PatchDraft, http.MethodPatch, "{}"},
		{"Solicitation.Submit", solHandler.Submit, http.MethodPost, ""},
		{"Solicitation.Approve", solHandler.Approve, http.MethodPost, ""},
		{"Solicitation.Reject", solHandler.Reject, http.MethodPost, "{}"},
		{"Solicitation.PatchAnalysis", solHandler.PatchAnalysis, http.MethodPatch, "{}"},
		{"Solicitation.Finalize", solHandler.Finalize, http.MethodPost, "{}"},
		{"Solicitation.GetHistory", solHandler.GetHistory, http.MethodGet, ""},
		{"Category.Create", catHandler.Create, http.MethodPost, "{}"},
		{"Category.Update", catHandler.Update, http.MethodPatch, "{}"},
		{"Category.SetApprovers", catHandler.SetApprovers, http.MethodPatch, "{}"},
		{"User.List", userHandler.List, http.MethodGet, ""},
		{"User.Create", userHandler.Create, http.MethodPost, "{}"},
		{"Auth.Me", authHandler.Me, http.MethodGet, ""},
		{"Dashboard.Get", dashHandler.Get, http.MethodGet, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()

			tc.fn(rec, req)

			require.Equal(t, http.StatusUnauthorized, rec.Code)
		})
	}
}
