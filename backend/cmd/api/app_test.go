//go:build integration

package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/config"
)

// TestBuildApp_WiresRealDependencies sobe a aplicação inteira (banco real,
// migrations, seed, todos os serviços e handlers) exatamente como main()
// faz, e confere que o resultado responde de ponta a ponta — a mesma
// montagem que roda em produção, só sem o ListenAndServe bloqueante.
func TestBuildApp_WiresRealDependencies(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL não definida — pulando teste de integração")
	}

	cfg := &config.Config{
		DatabaseURL: dsn,
		JWTSecret:   "test-secret",
		CORSOrigin:  "http://localhost:5173",
		Port:        "0",
		Env:         "test",
	}

	router, pool, err := buildApp(context.Background(), cfg)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	srv := httptest.NewServer(router)
	t.Cleanup(srv.Close)

	resp, err := http.Get(srv.URL + "/healthz")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp2, err := http.Get(srv.URL + "/api/v1/auth/profiles")
	require.NoError(t, err)
	defer resp2.Body.Close()
	require.Equal(t, http.StatusOK, resp2.StatusCode)
}

func TestBuildApp_InvalidDatabaseURL(t *testing.T) {
	// Porta local sem nada escutando: a conexão é recusada de imediato, sem
	// depender de timeout de DNS.
	cfg := &config.Config{DatabaseURL: "postgres://postgres:postgres@127.0.0.1:59999/nonexistent?sslmode=disable&connect_timeout=2", Port: "0"}

	_, _, err := buildApp(context.Background(), cfg)

	require.Error(t, err)
}
