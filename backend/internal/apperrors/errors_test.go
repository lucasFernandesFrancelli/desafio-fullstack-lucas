package apperrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"ekaizen-backend/internal/apperrors"
)

func TestAppError_Error(t *testing.T) {
	err := apperrors.NotFound("não encontrado")
	assert.Equal(t, "não encontrado", err.Error())
}

func TestConstructors(t *testing.T) {
	cases := []struct {
		name string
		err  *apperrors.AppError
		code apperrors.Code
	}{
		{"NotFound", apperrors.NotFound("x"), apperrors.CodeNotFound},
		{"Forbidden", apperrors.Forbidden("x"), apperrors.CodeForbidden},
		{"Conflict", apperrors.Conflict("x"), apperrors.CodeConflict},
		{"Unauthorized", apperrors.Unauthorized("x"), apperrors.CodeUnauthorized},
		{"Validation", apperrors.Validation("x", map[string]string{"f": "obrigatório"}), apperrors.CodeValidation},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.code, c.err.Code)
		})
	}
}

func TestAs_MatchesAppError(t *testing.T) {
	err := apperrors.Conflict("conflito")
	got, ok := apperrors.As(err)
	assert.True(t, ok)
	assert.Equal(t, err, got)
}

func TestAs_RejectsPlainError(t *testing.T) {
	_, ok := apperrors.As(errors.New("erro genérico"))
	assert.False(t, ok)
}
