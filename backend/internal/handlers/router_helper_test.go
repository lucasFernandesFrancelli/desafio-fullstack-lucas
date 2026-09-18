// Testes de handler: montam o router real (server.NewRouter) com serviços
// reais, mas repositórios mockados — exercitam parsing de request, códigos
// de status HTTP e o formato do JSON de erro, sem precisar de banco.
package handlers_test

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/auth"
	"ekaizen-backend/internal/handlers"
	"ekaizen-backend/internal/models"
	"ekaizen-backend/internal/server"
	"ekaizen-backend/internal/services"
	"ekaizen-backend/internal/testutil"
)

type testDeps struct {
	router *httptest.Server
	sols   *testutil.SolicitationRepo
	cats   *testutil.CategoryRepo
	hist   *testutil.HistoryRepo
	users  *testutil.UserRepo
	issuer *auth.JWTIssuer
}

func newTestServer(t *testing.T) *testDeps {
	t.Helper()

	issuer := auth.NewJWTIssuer("test-secret")
	sols := &testutil.SolicitationRepo{}
	cats := &testutil.CategoryRepo{}
	hist := &testutil.HistoryRepo{}
	users := &testutil.UserRepo{}

	authSvc := services.NewAuthService(users, issuer)
	solSvc := services.NewSolicitationService(&testutil.Transactor{}, sols, cats, hist)
	dashSvc := services.NewDashboardService(sols)
	adminSvc := services.NewAdminService(&testutil.Transactor{}, users, cats)

	router := server.NewRouter(server.Dependencies{
		Issuer:              issuer,
		CORSOrigin:          "http://localhost:5173",
		AuthHandler:         handlers.NewAuthHandler(authSvc),
		CategoryHandler:     handlers.NewCategoryHandler(cats, adminSvc),
		SolicitationHandler: handlers.NewSolicitationHandler(solSvc),
		DashboardHandler:    handlers.NewDashboardHandler(dashSvc),
		UserHandler:         handlers.NewUserHandler(adminSvc),
	})

	srv := httptest.NewServer(router)
	t.Cleanup(srv.Close)

	return &testDeps{router: srv, sols: sols, cats: cats, hist: hist, users: users, issuer: issuer}
}

func (d *testDeps) bearerFor(t *testing.T, u models.User) string {
	t.Helper()
	token, err := d.issuer.Issue(u)
	require.NoError(t, err)
	return "Bearer " + token
}
