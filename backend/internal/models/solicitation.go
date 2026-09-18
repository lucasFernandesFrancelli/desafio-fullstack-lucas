package models

import (
	"time"

	"github.com/google/uuid"
)

// Solicitation é o agregado central do domínio. Campos que só fazem sentido
// em determinados estados (categoria, notas de análise, motivo de recusa)
// são ponteiros/nullable porque nascem vazios em "rascunho".
type Solicitation struct {
	ID                  uuid.UUID  `json:"id"`
	Title               string     `json:"title"`
	ProblemDescription  string     `json:"problemDescription"`
	ProposedImprovement string     `json:"proposedImprovement"`
	CategoryID          *uuid.UUID `json:"categoryId"`
	CategoryName        string     `json:"categoryName,omitempty"`
	Location            string     `json:"location"`
	RequesterID         uuid.UUID  `json:"requesterId"`
	RequesterName       string     `json:"requesterName,omitempty"`
	Status              Status     `json:"status"`
	CurrentApprovalStep *int16     `json:"currentApprovalStep"`
	Severity            *int16     `json:"severity"`
	Urgency             *int16     `json:"urgency"`
	Trend               *int16     `json:"trend"`
	Priority            *int16     `json:"priority"`
	AnalysisNotes       string     `json:"analysisNotes"`
	RejectionReason     string     `json:"rejectionReason"`
	LastTransitionAt    time.Time  `json:"lastTransitionAt"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
}

// PendingActor descreve quem precisa agir agora numa solicitação. Pode ser
// uma pessoa específica (aprovador da vez, dono do rascunho) ou um papel
// (qualquer analista), ou ninguém quando o estado é terminal.
type PendingActor struct {
	Kind string  `json:"kind"` // "user" | "role" | "none"
	ID   *string `json:"id,omitempty"`
	Name string  `json:"name"`
}

// Permissions é calculado inteiramente no backend a partir do usuário
// autenticado e do estado atual — o frontend só espelha essas flags.
type Permissions struct {
	CanEdit     bool `json:"canEdit"`
	CanSubmit   bool `json:"canSubmit"`
	CanApprove  bool `json:"canApprove"`
	CanReject   bool `json:"canReject"`
	CanAnalyze  bool `json:"canAnalyze"`
	CanFinalize bool `json:"canFinalize"`
}

type HistoryEntry struct {
	ID             uuid.UUID     `json:"id"`
	SolicitationID uuid.UUID     `json:"solicitationId"`
	FromStatus     *Status       `json:"fromStatus"`
	ToStatus       Status        `json:"toStatus"`
	Action         HistoryAction `json:"action"`
	ActorID        uuid.UUID     `json:"actorId"`
	ActorName      string        `json:"actorName"`
	ApprovalStep   *int16        `json:"approvalStep,omitempty"`
	Comment        string        `json:"comment,omitempty"`
	Severity       *int16        `json:"severity,omitempty"`
	Urgency        *int16        `json:"urgency,omitempty"`
	Trend          *int16        `json:"trend,omitempty"`
	CreatedAt      time.Time     `json:"createdAt"`
}

// SolicitationDetail agrega tudo que a tela de detalhe precisa numa única
// resposta: os dados, quem precisa agir, o que o usuário atual pode fazer,
// e o histórico completo.
type SolicitationDetail struct {
	Solicitation
	PendingActor PendingActor   `json:"pendingActor"`
	Permissions  Permissions    `json:"permissions"`
	History      []HistoryEntry `json:"history"`
}

// SolicitationSummary é o formato enxuto usado por Lista e Kanban — ambas
// as visões consomem exatamente os mesmos dados de GET /solicitations.
type SolicitationSummary struct {
	ID                  uuid.UUID `json:"id"`
	Title               string    `json:"title"`
	CategoryName        string    `json:"categoryName,omitempty"`
	Status              Status    `json:"status"`
	CurrentApprovalStep *int16    `json:"currentApprovalStep"`
	RequesterName       string    `json:"requesterName"`
	Priority            *int16    `json:"priority"`
	PendingActorName    string    `json:"pendingActorName"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

// SolicitationFilter representa os parâmetros de busca/filtro usados tanto
// pela Lista quanto pelo Kanban (mesma fonte de dados).
type SolicitationFilter struct {
	Status     *Status
	CategoryID *uuid.UUID
	Query      string
}

// DraftInput carrega os campos editáveis de um rascunho (criação ou PATCH
// parcial) — todos opcionais porque o rascunho pode ser salvo incompleto.
type DraftInput struct {
	Title               *string    `json:"title"`
	ProblemDescription  *string    `json:"problemDescription"`
	ProposedImprovement *string    `json:"proposedImprovement"`
	CategoryID          *uuid.UUID `json:"categoryId"`
	Location            *string    `json:"location"`
}

// AnalysisInput carrega os campos do parecer, usados tanto no PATCH parcial
// quanto (todos preenchidos) na finalização.
type AnalysisInput struct {
	Severity      *int16  `json:"severity"`
	Urgency       *int16  `json:"urgency"`
	Trend         *int16  `json:"trend"`
	AnalysisNotes *string `json:"analysisNotes"`
}
