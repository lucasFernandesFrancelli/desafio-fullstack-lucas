package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"ekaizen-backend/internal/models"
)

func TestRole_Valid(t *testing.T) {
	assert.True(t, models.RoleColaborador.Valid())
	assert.True(t, models.RoleAnalista.Valid())
	assert.True(t, models.RoleGestor.Valid())
	assert.False(t, models.Role("inexistente").Valid())
}

func TestStatus_Terminal(t *testing.T) {
	assert.True(t, models.StatusFinalizada.Terminal())
	assert.True(t, models.StatusRecusada.Terminal())
	assert.False(t, models.StatusRascunho.Terminal())
	assert.False(t, models.StatusEmAprovacao.Terminal())
	assert.False(t, models.StatusEmAnalise.Terminal())
}
