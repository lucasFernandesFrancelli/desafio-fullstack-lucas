// Package apperrors define os erros de domínio usados pelos serviços. Os
// handlers HTTP mapeiam cada um para o status code correspondente, mantendo
// as regras de negócio livres de detalhes de transporte.
package apperrors

import "errors"

type Code string

const (
	CodeNotFound     Code = "not_found"
	CodeForbidden    Code = "forbidden"
	CodeConflict     Code = "conflict"
	CodeValidation   Code = "validation_error"
	CodeUnauthorized Code = "unauthorized"
)

type AppError struct {
	Code    Code
	Message string
	Fields  map[string]string
}

func (e *AppError) Error() string {
	return e.Message
}

func NotFound(message string) *AppError {
	return &AppError{Code: CodeNotFound, Message: message}
}

func Forbidden(message string) *AppError {
	return &AppError{Code: CodeForbidden, Message: message}
}

func Conflict(message string) *AppError {
	return &AppError{Code: CodeConflict, Message: message}
}

func Unauthorized(message string) *AppError {
	return &AppError{Code: CodeUnauthorized, Message: message}
}

func Validation(message string, fields map[string]string) *AppError {
	return &AppError{Code: CodeValidation, Message: message, Fields: fields}
}

// As converte err para *AppError se possível, retornando ok=false caso
// contrário (erro inesperado, deve virar 500).
func As(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
