# Da ideia à decisão — eKaizen

Sistema de gestão de sugestões de melhoria operacional: do registro da solicitação até a decisão final, passando por aprovação dupla sequencial e análise técnica. Desenvolvido como resposta ao desafio técnico "Da ideia à decisão" (eKaizen).

**Links da entrega**
- Aplicação (Vercel): https://desafio-fullstack-lucas.vercel.app
- API (Render): https://desafio-fullstack-lucas.onrender.com
- Repositório: https://github.com/lucasFernandesFrancelli/desafio-fullstack-lucas

## Sumário

- [Visão geral](#visão-geral)
- [Stack](#stack)
- [Arquitetura e decisões de design](#arquitetura-e-decisões-de-design)
- [Como rodar localmente](#como-rodar-localmente)
- [Usuários de teste (seed)](#usuários-de-teste-seed)
- [Testes e cobertura](#testes-e-cobertura)
- [Deploy](#deploy)
- [Uso de IA no desenvolvimento](#uso-de-ia-no-desenvolvimento)

## Visão geral

Uma empresa recebe sugestões de melhoria de quem vive a operação. A aplicação resolve essa rotina do início ao fim:

1. **Solicitante** cria uma solicitação (título, descrição do problema, melhoria proposta, categoria, local). Pode salvar como **rascunho** incompleto e continuar depois.
2. Ao enviar, a solicitação entra em **aprovação**: cada categoria tem dois aprovadores fixos, em ordem — só o aprovador da vez pode decidir. Qualquer um dos dois pode **recusar** (com motivo obrigatório), encerrando o fluxo. A segunda aprovação libera a etapa seguinte.
3. Em **análise**, um analista registra um parecer e notas de 1 a 5 para gravidade, urgência e tendência. A prioridade é o produto das três (1 a 125). A finalização exige todos esses dados.
4. O item chega a **finalizada** ou **recusada** — estados terminais, disponíveis para consulta.
5. **Lista** e **Kanban** mostram as mesmas solicitações, com busca e filtros por estado/categoria. Abrir uma solicitação mostra contexto, etapa atual, quem precisa agir e o histórico completo de decisões.
6. Um **dashboard** dá ao gestor uma visão agregada: quantas solicitações há em cada estado, quais estão paradas há mais tempo, e quem tem mais pendências — sem precisar perguntar pessoa por pessoa.
7. Uma área de **Gestão** (menu "Gestão", só para o gestor) mostra quem são os dois aprovadores de cada categoria com dados reais do banco — e permite criar categorias, reatribuir aprovadores e cadastrar novas pessoas, para confirmar que a regra não está fixa no código.

### Decisões de design assumidas

- **Login**: não há senha. Uma tela de seletor de perfil lista os usuários fictícios e um clique loga como aquele usuário — o backend ainda emite um JWT real, só não exige senha digitada (ideal para um avaliador alternar entre papéis rapidamente).
- **Visibilidade de rascunhos**: só o dono e o gestor veem uma solicitação em rascunho, em Lista/Kanban/Detalhe. Os demais estados são visíveis a todos os perfis autenticados; as ações (aprovar, recusar, analisar) ficam restritas por permissão, calculada inteiramente no backend e apenas espelhada pelo frontend.
- **Kanban**: somente visualização — cards não são arrastáveis. Todas as ações acontecem na tela de detalhe, aberta ao clicar no card (tanto na Lista quanto no Kanban).
- **Categoria como ligação, não papel fixo**: qualquer usuário pode criar solicitações. "Aprovador" é uma relação (categoria + ordem 1 ou 2) e "analista"/"gestor" são papéis (`role`) do usuário — refletindo a regra de que os dois aprovadores de uma categoria são pessoas específicas e pré-definidas.

## Stack

**Backend** (`/backend`): Go 1.23+, [chi](https://github.com/go-chi/chi) (router), [pgx v5](https://github.com/jackc/pgx) (driver Postgres, sem ORM — repositório com interfaces + SQL manual), JWT (HS256) para o "login" mock, migrations embutidas no binário (`embed.FS`), seed idempotente que roda no boot.

**Frontend** (`/frontend`): React 19, TypeScript em modo estrito, Vite, Tailwind CSS 4, [shadcn/ui](https://ui.shadcn.com), TanStack Query, React Router, react-hook-form + zod.

**Testes**: Go `testing` + `testify` (unitário + integração contra Postgres real) · Vitest + React Testing Library + MSW (unitário/componente) · Playwright (ponta a ponta, contra a API e o Postgres reais).

**Banco**: PostgreSQL (local via serviço nativo em desenvolvimento; [Neon](https://neon.tech) serverless em produção).

## Arquitetura e decisões de design

### Backend

`SolicitationService` (`backend/internal/services/solicitation_service.go`) é a **única porta de entrada** para criar ou mudar o estado de uma solicitação. Handlers HTTP nunca tocam o repositório diretamente — isso garante que Lista, Kanban e Detalhe, todos consumindo `GET /solicitations` e `GET /solicitations/{id}`, respeitem exatamente as mesmas regras. Cada transição roda dentro de uma transação com `SELECT ... FOR UPDATE`, evitando corrida em cliques duplos.

Máquina de estados:

| Transição | De → Para | Quem pode | Campos obrigatórios |
|---|---|---|---|
| Criar | — → `rascunho` | qualquer autenticado | título |
| Enviar | `rascunho` → `em_aprovacao` (etapa 1) | dono | todos os campos |
| Aprovar (1ª) | `em_aprovacao` etapa 1 → etapa 2 | aprovador de ordem 1 da categoria | — |
| Aprovar (2ª) | `em_aprovacao` etapa 2 → `em_analise` | aprovador de ordem 2 da categoria | — |
| Recusar | `em_aprovacao` → `recusada` | aprovador da vez | motivo |
| Finalizar | `em_analise` → `finalizada` | qualquer `analista` | parecer + gravidade + urgência + tendência |

`finalizada` e `recusada` são terminais. O detalhe de cada solicitação expõe `pendingActor` (quem precisa agir) e `permissions` (o que o usuário logado pode fazer ali) — ambos calculados no servidor.

Sem ORM: `internal/repository/postgres` escreve SQL diretamente contra `pgx`, atrás de interfaces (`internal/repository/interfaces.go`) que são mockadas nos testes de serviço com `testify/mock`.

### Frontend

Nenhum componente reimplementa RBAC: o objeto `permissions` retornado pela API decide o que cada painel de ação (`ApproveRejectPanel`, `AnalysisPanel`, edição do formulário) mostra. `TanStack Query` mantém a Lista, o Kanban e o Detalhe sincronizados — qualquer mutação invalida as mesmas chaves de cache.

## Como rodar localmente

### Pré-requisitos

- [Go 1.23+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org)
- PostgreSQL 16+ (local ou um banco na nuvem, ex. Neon) — **Docker não é necessário**, o projeto não depende dele.

### 1. Banco de dados

Crie dois bancos vazios (dev e teste):

```bash
psql -U postgres -c "CREATE DATABASE ekaizen_dev;"
psql -U postgres -c "CREATE DATABASE ekaizen_test;"
```

### 2. Backend

```bash
cd backend
cp .env.example .env
# edite .env se sua senha/porta do Postgres for diferente de postgres/5432
go run ./cmd/api
```

Migrations e seed rodam automaticamente no boot. A API sobe em `http://localhost:8080`, com `/healthz` para checagem e `/api/v1/*` para as rotas. Para reaplicar o seed manualmente (idempotente): `go run ./cmd/seed`.

### 3. Frontend

Em outro terminal:

```bash
cd frontend
cp .env.example .env.local
npm install
npm run dev
```

Abra `http://localhost:5173`. A tela inicial lista os usuários de teste — clique em qualquer um para entrar.

## Usuários de teste (seed)

Dados fictícios, recriados a cada boot do backend (idempotente). Categorias e seus dois aprovadores (ordem 1 → 2):

| Categoria | 1º aprovador | 2º aprovador |
|---|---|---|
| Segurança do Trabalho | Ricardo Nogueira | Fernanda Albuquerque |
| Qualidade | Camila Duarte | Marcos Vieira |
| Produtividade | Bruno Tanaka | Patrícia Lemos |
| Meio Ambiente | Diego Farias | Juliana Prado |

Outros perfis: **Renata Souza** e **Tiago Martins** (analistas, atendem qualquer categoria), **Sérgio Andrade** (gestor, acessa o dashboard e a área de gestão), **Ana Beatriz Costa**, **João Pedro Rocha** e **Larissa Mendes** (colaboradores, só solicitam).

O seed também cria uma solicitação de exemplo em cada um dos 5 estados, para o avaliador ver o processo completo sem precisar operá-lo manualmente do zero.

Logado como **Sérgio Andrade** (gestor), o menu **Gestão** mostra essa mesma tabela de aprovadores vinda do banco (não fixa no código) e permite criar categorias, reatribuir aprovadores ou cadastrar novas pessoas para testar a regra com dados novos.

## Testes e cobertura

### Backend (meta ≥80% — atingido: **90.6%** de cobertura de linhas)

```bash
cd backend

# Unitários (serviços mockados + handlers httptest) — não precisam de banco,
# pois os arquivos de integração ficam fora do build sem a tag abaixo
go test ./...

# Suíte completa (inclui integração de repositório/migrations/seed contra Postgres real) + relatório
export TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5432/ekaizen_test?sslmode=disable"
go test ./... -tags=integration -p 1 -coverpkg=./... -coverprofile=coverage.out -covermode=atomic
go tool cover -func=coverage.out    # resumo no terminal
go tool cover -html=coverage.out -o coverage.html   # relatório navegável
```

> `-p 1` evita que pacotes rodem em paralelo disputando as mesmas tabelas do banco de teste. `TEST_DATABASE_URL` ausente faz os testes de integração pularem automaticamente (`t.Skip`), então `go test ./...` sozinho sempre funciona, mesmo sem Postgres disponível.
>
> No Windows, rode o comando de cobertura acima via um shell **bash** (Git Bash/WSL), não PowerShell: em uma suíte de múltiplos pacotes o PowerShell corrompeu silenciosamente o merge do `coverage.out` (testes passando normalmente, mas os blocos cobertos por testes de branches de erro apareciam como não exercitados no relatório final). Rodando o mesmo comando via bash o relatório bate com a execução real.

Cobertura por camada: serviços (regras de negócio) e handlers via mocks/`httptest`; repositórios via integration tests contra Postgres real (valida inclusive constraints do schema, como os dois aprovadores distintos por categoria). `cmd/api`/`cmd/seed` foram fatorados em `buildApp`/`runSeed` justamente para serem exercitados por um teste de integração de ponta a ponta, deixando só a função `main()` (glue de bootstrap) fora da meta.

Notas sobre como a cobertura foi de **80,6% para 90,6%** (branches de erro que um teste "caminho feliz" nunca alcança):

- **Actor ausente**: quase todo handler tem um `if !ok { return 401 }` defensivo para quando o middleware não injeta o ator no contexto — inalcançável passando pelo router normal, então esses testes chamam o método do handler diretamente, sem passar pelo middleware.
- **Erros genéricos de repositório**: para cobrir o `return err` "não é o erro específico esperado" (ex.: violação de unicidade vs. qualquer outro erro do Postgres), os testes passam um `context.Context` já cancelado para o repositório — o pgx recusa a query e devolve um erro genérico, sem precisar quebrar o schema do banco.
- **Migration "do zero"**: como o mesmo Postgres de teste é reaproveitado por vários pacotes, a migration já estava marcada como aplicada bem antes do teste de `migrations.Apply` rodar — o teste só exercitava o branch de idempotência. A correção foi resetar o schema (`DROP SCHEMA public CASCADE`) no início do teste, forçando o `Apply` "de verdade" a rodar.
- **Erros do seed**: para cobrir os `fmt.Errorf(...)` de cada etapa do `seed.Run` (usuários, categorias, aprovadores, solicitações de exemplo), os testes resetam o schema, aplicam as migrations e então derrubam (`DROP TABLE ... CASCADE`) a tabela específica que aquela etapa precisa, forçando o erro sem mexer em permissões (o usuário do Postgres de teste normalmente é superuser, então `REVOKE` não teria efeito).

### Frontend (meta ≥80% — atingido: **90%** de statements / **92%** de funções)

```bash
cd frontend

npm run test              # roda a suíte (Vitest + Testing Library + MSW)
npm run test:coverage     # com relatório de cobertura (texto + HTML em frontend/coverage)
```

### Ponta a ponta (Playwright)

6 especificações cobrindo o desafio inteiro contra a API e o Postgres reais (sobe automaticamente uma instância da API em `:8081` e do Vite em `:5174`, apontando para um banco `ekaizen_e2e` dedicado):

```bash
cd frontend
createdb ekaizen_e2e   # ou: psql -U postgres -c "CREATE DATABASE ekaizen_e2e;"
npx playwright install chromium   # primeira vez apenas
npm run e2e
```

- `golden-path.spec.ts` — rascunho → editar → enviar → aprovar (1ª) → aprovar (2ª) → parecer + notas → finalizar; confere a prioridade calculada e a paridade Lista/Kanban.
- `rejection-path.spec.ts` — recusa com motivo, estado terminal, visível para o gestor.
- `draft-validation.spec.ts` — bloqueio client-side ao tentar salvar/enviar incompleto.
- `rbac-guard.spec.ts` — aprovador fora da vez e analista não veem as ações; uma chamada direta à API é recusada com 403.
- `kanban-list-parity.spec.ts` — mesmo filtro retorna o mesmo conjunto nas duas visões; o card do Kanban não é arrastável.
- `admin-management.spec.ts` — gestor cria usuário e categoria (com os dois aprovadores), reatribui um aprovador, e confirma que um colaborador não acessa `/admin`.

## Deploy

Arquitetura de produção: **Vercel** (frontend) → **Render** (API Go, Docker) → **Neon** (Postgres serverless). CORS usa apenas Bearer token (sem cookies), evitando a complexidade de `SameSite`/`credentials` entre domínios diferentes.

1. **Neon**: criar um projeto, copiar a `DATABASE_URL` (com `sslmode=require`).
2. **Render**: novo Web Service, "Docker", root do repo `backend/` (usa `backend/render.yaml` e `backend/Dockerfile`). Variáveis: `DATABASE_URL` (a do Neon), `JWT_SECRET` (gerado automaticamente), `CORS_ORIGIN` (preenchida depois de criar o projeto na Vercel). Migrations e seed rodam sozinhos no boot do container.
3. **Vercel**: novo projeto apontando para este repositório com **Root Directory = `frontend`**. Variável `VITE_API_BASE_URL = https://<seu-app>.onrender.com/api/v1`.
4. Voltar ao Render e atualizar `CORS_ORIGIN` com a URL final da Vercel; redeploy.
5. Testar o fluxo completo nas URLs públicas.

> **Nota**: o plano gratuito do Render "dorme" a API após períodos de inatividade — a primeira requisição depois de um tempo sem uso pode levar de 30 a 60 segundos para responder enquanto o container acorda. As chamadas seguintes voltam ao normal.

## Uso de IA no desenvolvimento

Todo o projeto foi desenvolvido com [Claude Code](https://claude.com/claude-code) (Anthropic), da arquitetura e do modelo de dados até a implementação, os testes e este README — atendendo ao requisito do desafio de usar IA no desenvolvimento.
