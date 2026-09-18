// Command api sobe o servidor HTTP: carrega config, conecta no Postgres,
// aplica migrations e seed, monta as dependências e escuta.
package main

import (
	"context"
	"log"
	"net/http"

	"ekaizen-backend/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	router, pool, err := buildApp(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	addr := ":" + cfg.Port
	log.Printf("ekaizen-backend ouvindo em %s (env=%s)", addr, cfg.Env)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("servidor: %v", err)
	}
}
