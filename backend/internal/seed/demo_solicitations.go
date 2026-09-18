package seed

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"ekaizen-backend/internal/models"
)

type demoSolicitation struct {
	key                 string
	title               string
	problem             string
	improvement         string
	categoryKey         string
	location            string
	requesterKey        string
	status              models.Status
	currentApprovalStep *int16
	severity            *int16
	urgency             *int16
	trend               *int16
	analysisNotes       string
	rejectionReason     string
}

func i16(v int16) *int16 { return &v }

// approverKeyFor espelha a configuração de aprovadores das categorias acima,
// usado só para montar um histórico de exemplo plausível.
func approverKeyFor(categoryKey string, order int16) string {
	for _, c := range categories {
		if c.key == categoryKey {
			if order == 1 {
				return c.approver1
			}
			return c.approver2
		}
	}
	return ""
}

var demoSolicitations = []demoSolicitation{
	{
		key:          "demo-rascunho",
		title:        "Sinalização insuficiente no pátio de cargas",
		problem:      "Empilhadeiras cruzam com pedestres sem sinalização no pátio de cargas.",
		improvement:  "Pintar faixas de pedestre e instalar espelhos convexos nos cruzamentos.",
		categoryKey:  "seguranca",
		location:     "Pátio de Cargas - Galpão 2",
		requesterKey: "ana",
		status:       models.StatusRascunho,
	},
	{
		key:                 "demo-aprovacao1",
		title:               "Retrabalho recorrente na linha de embalagem",
		problem:             "Itens saem com etiqueta trocada, gerando retrabalho diário.",
		improvement:         "Checklist de dupla checagem antes do lacre da caixa.",
		categoryKey:         "qualidade",
		location:            "Linha de Embalagem 3",
		requesterKey:        "joao",
		status:              models.StatusEmAprovacao,
		currentApprovalStep: i16(1),
	},
	{
		key:                 "demo-aprovacao2",
		title:               "EPI incorreto para manuseio de produtos químicos",
		problem:             "As luvas atuais não protegem contra o solvente usado na limpeza industrial.",
		improvement:         "Trocar por luvas nitrílicas de cano longo conforme ficha de segurança.",
		categoryKey:         "seguranca",
		location:            "Setor de Limpeza Industrial",
		requesterKey:        "larissa",
		status:              models.StatusEmAprovacao,
		currentApprovalStep: i16(2),
	},
	{
		key:          "demo-analise",
		title:        "Tempo de setup elevado na troca de molde",
		problem:      "A troca de molde leva 90 minutos, parando a linha inteira.",
		improvement:  "Aplicar técnica SMED e preparar ferramentas fora do ciclo de parada.",
		categoryKey:  "produtividade",
		location:     "Prensa 4 - Setor de Estamparia",
		requesterKey: "joao",
		status:       models.StatusEmAnalise,
	},
	{
		key:           "demo-finalizada",
		title:         "Descarte incorreto de resíduos de tinta",
		problem:       "Resíduos de tinta vão para o lixo comum, contaminando o solo ao redor do pátio.",
		improvement:   "Instalar coletores específicos e contratar destinação certificada.",
		categoryKey:   "ambiente",
		location:      "Cabine de Pintura 1",
		requesterKey:  "ana",
		status:        models.StatusFinalizada,
		severity:      i16(4),
		urgency:       i16(3),
		trend:         i16(5),
		analysisNotes: "Risco ambiental confirmado; prioridade alta para tratativa em até 30 dias.",
	},
	{
		key:             "demo-recusada",
		title:           "Trocar todos os monitores por modelos curvos",
		problem:         "Colaboradores relatam preferência por monitores curvos por conforto visual.",
		improvement:     "Substituir os 40 monitores do escritório administrativo.",
		categoryKey:     "qualidade",
		location:        "Escritório Administrativo",
		requesterKey:    "larissa",
		status:          models.StatusRecusada,
		rejectionReason: "Fora do escopo de melhoria de processo; sugerido abrir como solicitação de facilities.",
	},
}

func seedDemoSolicitations(ctx context.Context, pool *pgxpool.Pool) error {
	for _, d := range demoSolicitations {
		solID := id("solicitation:" + d.key)

		var exists bool
		if err := pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM solicitations WHERE id = $1)`, solID).Scan(&exists); err != nil {
			return err
		}
		if exists {
			continue
		}

		if _, err := pool.Exec(ctx, `
			INSERT INTO solicitations
				(id, title, problem_description, proposed_improvement, category_id, location, requester_id,
				 status, current_approval_step, severity, urgency, trend, analysis_notes, rejection_reason)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
			solID, d.title, d.problem, d.improvement, id("category:"+d.categoryKey), d.location, id("user:"+d.requesterKey),
			d.status, d.currentApprovalStep, d.severity, d.urgency, d.trend, d.analysisNotes, d.rejectionReason,
		); err != nil {
			return err
		}

		if err := seedHistoryFor(ctx, pool, solID, d); err != nil {
			return err
		}
	}
	return nil
}

func seedHistoryFor(ctx context.Context, pool *pgxpool.Pool, solID uuid.UUID, d demoSolicitation) error {
	rascunho, emAprovacao, emAnalise := models.StatusRascunho, models.StatusEmAprovacao, models.StatusEmAnalise

	insert := func(suffix string, from *models.Status, to models.Status, action models.HistoryAction,
		actorKey string, approvalStep *int16, comment string, severity, urgency, trend *int16) error {
		_, err := pool.Exec(ctx, `
			INSERT INTO solicitation_history
				(id, solicitation_id, from_status, to_status, action, actor_id, approval_step, comment, severity, urgency, trend)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
			id("history:"+d.key+":"+suffix), solID, from, to, action, id("user:"+actorKey), approvalStep, comment, severity, urgency, trend)
		return err
	}

	if err := insert("criada", nil, models.StatusRascunho, models.ActionCriada, d.requesterKey, nil, "", nil, nil, nil); err != nil {
		return err
	}

	approver1, approver2 := approverKeyFor(d.categoryKey, 1), approverKeyFor(d.categoryKey, 2)

	switch d.key {
	case "demo-rascunho":
		// só o evento de criação mesmo — ainda não foi enviada.

	case "demo-aprovacao1":
		if err := insert("enviada", &rascunho, emAprovacao, models.ActionEnviada, d.requesterKey, nil, "", nil, nil, nil); err != nil {
			return err
		}

	case "demo-aprovacao2":
		if err := insert("enviada", &rascunho, emAprovacao, models.ActionEnviada, d.requesterKey, nil, "", nil, nil, nil); err != nil {
			return err
		}
		if err := insert("aprovada1", &emAprovacao, emAprovacao, models.ActionAprovadaEtapa1, approver1, i16(1), "", nil, nil, nil); err != nil {
			return err
		}

	case "demo-analise":
		if err := insert("enviada", &rascunho, emAprovacao, models.ActionEnviada, d.requesterKey, nil, "", nil, nil, nil); err != nil {
			return err
		}
		if err := insert("aprovada1", &emAprovacao, emAprovacao, models.ActionAprovadaEtapa1, approver1, i16(1), "", nil, nil, nil); err != nil {
			return err
		}
		if err := insert("aprovada2", &emAprovacao, emAnalise, models.ActionAprovadaEtapa2, approver2, i16(2), "", nil, nil, nil); err != nil {
			return err
		}

	case "demo-finalizada":
		if err := insert("enviada", &rascunho, emAprovacao, models.ActionEnviada, d.requesterKey, nil, "", nil, nil, nil); err != nil {
			return err
		}
		if err := insert("aprovada1", &emAprovacao, emAprovacao, models.ActionAprovadaEtapa1, approver1, i16(1), "", nil, nil, nil); err != nil {
			return err
		}
		if err := insert("aprovada2", &emAprovacao, emAnalise, models.ActionAprovadaEtapa2, approver2, i16(2), "", nil, nil, nil); err != nil {
			return err
		}
		if err := insert("finalizada", &emAnalise, models.StatusFinalizada, models.ActionFinalizada, "renata", nil,
			d.analysisNotes, d.severity, d.urgency, d.trend); err != nil {
			return err
		}

	case "demo-recusada":
		if err := insert("enviada", &rascunho, emAprovacao, models.ActionEnviada, d.requesterKey, nil, "", nil, nil, nil); err != nil {
			return err
		}
		if err := insert("recusada", &emAprovacao, models.StatusRecusada, models.ActionRecusada, approver1, i16(1),
			d.rejectionReason, nil, nil, nil); err != nil {
			return err
		}
	}

	return nil
}
