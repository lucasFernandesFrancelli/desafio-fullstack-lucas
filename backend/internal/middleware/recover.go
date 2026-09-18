package middleware

import (
	"fmt"
	"net/http"

	"ekaizen-backend/internal/httpx"
)

// Recover captura panics em handlers e responde 500 em vez de derrubar o
// processo, mantendo o formato de erro padrão da API.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				httpx.Error(w, fmt.Errorf("panic recuperado: %v", rec))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
