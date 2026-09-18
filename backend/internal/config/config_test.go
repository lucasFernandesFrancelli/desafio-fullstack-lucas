package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/config"
)

func TestLoad_Success(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://x")
	t.Setenv("PORT", "9090")
	t.Setenv("JWT_SECRET", "segredo")
	t.Setenv("CORS_ORIGIN", "https://example.com")
	t.Setenv("APP_ENV", "test")

	cfg, err := config.Load()

	require.NoError(t, err)
	assert.Equal(t, "postgres://x", cfg.DatabaseURL)
	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "segredo", cfg.JWTSecret)
	assert.Equal(t, "https://example.com", cfg.CORSOrigin)
	assert.Equal(t, "test", cfg.Env)
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	_, err := config.Load()

	assert.Error(t, err)
}

func TestLoad_UsesDefaultsWhenUnset(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://x")
	t.Setenv("PORT", "")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("CORS_ORIGIN", "")
	t.Setenv("APP_ENV", "")

	cfg, err := config.Load()

	require.NoError(t, err)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "dev-secret-change-me", cfg.JWTSecret)
	assert.Equal(t, "http://localhost:5173", cfg.CORSOrigin)
	assert.Equal(t, "development", cfg.Env)
}
