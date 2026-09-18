// Command seed aplica migrations + seed manualmente, útil para reset local
// (o servidor da API já faz isso sozinho em todo boot).
package main

import (
	"context"
	"log"

	"ekaizen-backend/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if err := runSeed(context.Background(), cfg.DatabaseURL); err != nil {
		log.Fatal(err)
	}

	log.Println("seed concluído com sucesso")
}
