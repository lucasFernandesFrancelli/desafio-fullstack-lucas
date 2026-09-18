// Package seed popula o banco com dados fictícios temáticos de melhoria
// contínua industrial: usuários, categorias com seus 2 aprovadores cada, e
// uma solicitação de exemplo em cada um dos 5 estados. Roda automaticamente
// no boot do backend (a avaliação é via link público, ninguém vai rodar um
// script manual) e é idempotente — IDs são determinísticos (uuid.NewSHA1)
// para que rodar de novo nunca duplique nada.
package seed

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"ekaizen-backend/internal/models"
)

var namespace = uuid.MustParse("2f6c8d0a-6b1a-4b8a-9c1a-3a1e9f6c8d0a")

func id(key string) uuid.UUID { return uuid.NewSHA1(namespace, []byte(key)) }

type seedUser struct {
	key   string
	name  string
	email string
	role  models.Role
}

type seedCategory struct {
	key         string
	name        string
	description string
	approver1   string
	approver2   string
}

var users = []seedUser{
	{"ricardo", "Ricardo Nogueira", "ricardo.nogueira@ekaizen.example", models.RoleColaborador},
	{"fernanda", "Fernanda Albuquerque", "fernanda.albuquerque@ekaizen.example", models.RoleColaborador},
	{"camila", "Camila Duarte", "camila.duarte@ekaizen.example", models.RoleColaborador},
	{"marcos", "Marcos Vieira", "marcos.vieira@ekaizen.example", models.RoleColaborador},
	{"bruno", "Bruno Tanaka", "bruno.tanaka@ekaizen.example", models.RoleColaborador},
	{"patricia", "Patrícia Lemos", "patricia.lemos@ekaizen.example", models.RoleColaborador},
	{"diego", "Diego Farias", "diego.farias@ekaizen.example", models.RoleColaborador},
	{"juliana", "Juliana Prado", "juliana.prado@ekaizen.example", models.RoleColaborador},
	{"renata", "Renata Souza", "renata.souza@ekaizen.example", models.RoleAnalista},
	{"tiago", "Tiago Martins", "tiago.martins@ekaizen.example", models.RoleAnalista},
	{"sergio", "Sérgio Andrade", "sergio.andrade@ekaizen.example", models.RoleGestor},
	{"ana", "Ana Beatriz Costa", "ana.costa@ekaizen.example", models.RoleColaborador},
	{"joao", "João Pedro Rocha", "joao.rocha@ekaizen.example", models.RoleColaborador},
	{"larissa", "Larissa Mendes", "larissa.mendes@ekaizen.example", models.RoleColaborador},
}

var categories = []seedCategory{
	{"seguranca", "Segurança do Trabalho", "Riscos e melhorias de segurança na operação.", "ricardo", "fernanda"},
	{"qualidade", "Qualidade", "Não conformidades e melhorias de processo.", "camila", "marcos"},
	{"produtividade", "Produtividade", "Eficiência e otimização de rotinas.", "bruno", "patricia"},
	{"ambiente", "Meio Ambiente", "Impactos ambientais e sustentabilidade.", "diego", "juliana"},
}

// Run aplica o seed. É seguro chamar em todo boot do processo.
func Run(ctx context.Context, pool *pgxpool.Pool) error {
	for _, u := range users {
		if _, err := pool.Exec(ctx, `
			INSERT INTO users (id, name, email, role) VALUES ($1,$2,$3,$4)
			ON CONFLICT (email) DO NOTHING`,
			id("user:"+u.key), u.name, u.email, u.role); err != nil {
			return fmt.Errorf("seed usuário %s: %w", u.key, err)
		}
	}

	for _, c := range categories {
		if _, err := pool.Exec(ctx, `
			INSERT INTO categories (id, name, description) VALUES ($1,$2,$3)
			ON CONFLICT (name) DO NOTHING`,
			id("category:"+c.key), c.name, c.description); err != nil {
			return fmt.Errorf("seed categoria %s: %w", c.key, err)
		}

		approvers := []struct {
			userKey string
			order   int16
		}{
			{c.approver1, 1},
			{c.approver2, 2},
		}
		for _, a := range approvers {
			approverID := id(fmt.Sprintf("approver:%s:%d", c.key, a.order))
			if _, err := pool.Exec(ctx, `
				INSERT INTO category_approvers (id, category_id, user_id, approval_order)
				VALUES ($1,$2,$3,$4)
				ON CONFLICT (category_id, approval_order) DO NOTHING`,
				approverID, id("category:"+c.key), id("user:"+a.userKey), a.order); err != nil {
				return fmt.Errorf("seed aprovador %s/%d: %w", c.key, a.order, err)
			}
		}
	}

	if err := seedDemoSolicitations(ctx, pool); err != nil {
		return fmt.Errorf("seed solicitações de exemplo: %w", err)
	}

	log.Println("seed: dados fictícios verificados/aplicados")
	return nil
}
