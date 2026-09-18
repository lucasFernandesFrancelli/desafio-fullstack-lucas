.PHONY: run seed test test-integration cover migrate-down

run:
	go run ./cmd/api

seed:
	go run ./cmd/seed

# Testes unitários (serviços mockados + handlers httptest) — não precisam de banco.
test:
	go test ./... -short -coverprofile=coverage.out -covermode=atomic

# Testes de integração de repositório — precisam de TEST_DATABASE_URL apontando
# para um Postgres real (local ou CI). -p 1 evita que pacotes rodem em paralelo
# e disputem/truquem as mesmas tabelas ao mesmo tempo.
test-integration:
	go test ./... -tags=integration -p 1 -coverprofile=coverage.out -covermode=atomic

cover: test-integration
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

migrate-down:
	psql "$$DATABASE_URL" -f migrations/000001_init.down.sql
