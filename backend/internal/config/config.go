// Package config carrega a configuração do serviço a partir de variáveis de
// ambiente (com .env opcional para desenvolvimento local).
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	CORSOrigin  string
	Env         string
}

// Load lê o .env (se existir) e depois as variáveis de ambiente reais, que
// sempre têm prioridade — necessário porque em produção (Render) não há
// arquivo .env, só env vars configuradas na plataforma.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-me"),
		CORSOrigin:  getEnv("CORS_ORIGIN", "http://localhost:5173"),
		Env:         getEnv("APP_ENV", "development"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL não configurada")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
