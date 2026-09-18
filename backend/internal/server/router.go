// Package server monta o roteador HTTP: quais rotas existem, quais exigem
// autenticação, e em que ordem os middlewares rodam.
package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"ekaizen-backend/internal/auth"
	"ekaizen-backend/internal/handlers"
	appmw "ekaizen-backend/internal/middleware"
)

type Dependencies struct {
	Issuer              *auth.JWTIssuer
	CORSOrigin          string
	AuthHandler         *handlers.AuthHandler
	CategoryHandler     *handlers.CategoryHandler
	SolicitationHandler *handlers.SolicitationHandler
	DashboardHandler    *handlers.DashboardHandler
	UserHandler         *handlers.UserHandler
}

func NewRouter(deps Dependencies) http.Handler {
	r := chi.NewRouter()
	r.Use(appmw.Recover)
	r.Use(appmw.Logging)
	r.Use(appmw.CORS(deps.CORSOrigin))

	r.Get("/healthz", handlers.Health)

	r.Route("/api/v1", func(r chi.Router) {
		// Rotas públicas: a tela de seletor de perfil precisa listar os
		// usuários e logar antes de ter qualquer token.
		r.Get("/auth/profiles", deps.AuthHandler.ListProfiles)
		r.Post("/auth/login", deps.AuthHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(appmw.Auth(deps.Issuer))

			r.Get("/auth/me", deps.AuthHandler.Me)
			r.Get("/dashboard", deps.DashboardHandler.Get)

			r.Route("/categories", func(r chi.Router) {
				r.Get("/", deps.CategoryHandler.List)
				r.Post("/", deps.CategoryHandler.Create)
				r.Patch("/{id}", deps.CategoryHandler.Update)
				r.Patch("/{id}/approvers", deps.CategoryHandler.SetApprovers)
			})

			r.Route("/users", func(r chi.Router) {
				r.Get("/", deps.UserHandler.List)
				r.Post("/", deps.UserHandler.Create)
			})

			r.Route("/solicitations", func(r chi.Router) {
				r.Post("/", deps.SolicitationHandler.Create)
				r.Get("/", deps.SolicitationHandler.List)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", deps.SolicitationHandler.Get)
					r.Patch("/", deps.SolicitationHandler.PatchDraft)
					r.Post("/submit", deps.SolicitationHandler.Submit)
					r.Post("/approve", deps.SolicitationHandler.Approve)
					r.Post("/reject", deps.SolicitationHandler.Reject)
					r.Patch("/analysis", deps.SolicitationHandler.PatchAnalysis)
					r.Post("/finalize", deps.SolicitationHandler.Finalize)
					r.Get("/history", deps.SolicitationHandler.GetHistory)
				})
			})
		})
	})

	return r
}
