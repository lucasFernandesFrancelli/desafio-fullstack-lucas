package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

// CORS libera o front local (Vite) e a origem de produção configurada via
// CORS_ORIGIN (a URL da Vercel). Usa apenas Bearer token, nunca cookies, o
// que evita toda a complexidade de SameSite/credentials entre domínios.
func CORS(productionOrigin string) func(http.Handler) http.Handler {
	origins := []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	if productionOrigin != "" {
		origins = append(origins, productionOrigin)
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}
