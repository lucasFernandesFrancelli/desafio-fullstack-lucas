package middleware

import (
	"context"
	"net/http"
	"strings"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/auth"
	"ekaizen-backend/internal/httpx"
	"ekaizen-backend/internal/models"
)

type ctxKey string

const actorCtxKey ctxKey = "actor"

// Auth exige um Bearer token válido e injeta o usuário autenticado no
// contexto da requisição para os handlers consumirem via ActorFrom.
func Auth(issuer *auth.JWTIssuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			token, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || token == "" {
				httpx.Error(w, apperrors.Unauthorized("token de autenticação ausente"))
				return
			}

			claims, err := issuer.Parse(token)
			if err != nil {
				httpx.Error(w, apperrors.Unauthorized("token inválido ou expirado"))
				return
			}

			actor := models.User{ID: claims.UserID, Name: claims.Name, Role: claims.Role}
			ctx := context.WithValue(r.Context(), actorCtxKey, actor)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ActorFrom(ctx context.Context) (models.User, bool) {
	actor, ok := ctx.Value(actorCtxKey).(models.User)
	return actor, ok
}
