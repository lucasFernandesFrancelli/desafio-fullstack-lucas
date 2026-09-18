// Package httpx contém helpers de resposta HTTP compartilhados por
// handlers e middlewares (mantido separado dos dois para evitar import
// cycle entre eles).
package httpx

import (
	"encoding/json"
	"log"
	"net/http"

	"ekaizen-backend/internal/apperrors"
)

func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// Error mapeia um erro de domínio (*apperrors.AppError) para o status HTTP
// correto; qualquer outro erro é tratado como falha interna (500) e logado,
// sem vazar detalhes internos na resposta.
func Error(w http.ResponseWriter, err error) {
	if appErr, ok := apperrors.As(err); ok {
		JSON(w, statusFor(appErr.Code), errorBody{Error: errorDetail{
			Code: string(appErr.Code), Message: appErr.Message, Fields: appErr.Fields,
		}})
		return
	}
	log.Printf("erro interno: %v", err)
	JSON(w, http.StatusInternalServerError, errorBody{Error: errorDetail{
		Code: "internal_error", Message: "erro interno inesperado",
	}})
}

func statusFor(code apperrors.Code) int {
	switch code {
	case apperrors.CodeNotFound:
		return http.StatusNotFound
	case apperrors.CodeForbidden:
		return http.StatusForbidden
	case apperrors.CodeConflict:
		return http.StatusConflict
	case apperrors.CodeValidation:
		return http.StatusUnprocessableEntity
	case apperrors.CodeUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
