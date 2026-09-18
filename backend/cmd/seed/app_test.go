//go:build integration

package main

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunSeed_Success(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL não definida — pulando teste de integração")
	}

	require.NoError(t, runSeed(context.Background(), dsn))
}

func TestRunSeed_InvalidDatabaseURL(t *testing.T) {
	err := runSeed(context.Background(), "postgres://postgres:postgres@127.0.0.1:59999/nonexistent?sslmode=disable&connect_timeout=2")
	require.Error(t, err)
}
