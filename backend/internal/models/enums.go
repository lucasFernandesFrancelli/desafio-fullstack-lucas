package models

// Role identifica o papel base de um usuário. Qualquer usuário pode criar
// solicitações independente do papel; role só concede capacidades extras
// (analista pode dar parecer, gestor enxerga o dashboard).
type Role string

const (
	RoleColaborador Role = "colaborador"
	RoleAnalista    Role = "analista"
	RoleGestor      Role = "gestor"
)

func (r Role) Valid() bool {
	switch r {
	case RoleColaborador, RoleAnalista, RoleGestor:
		return true
	}
	return false
}

// Status representa os 5 estados do processo definidos pelo desafio.
type Status string

const (
	StatusRascunho    Status = "rascunho"
	StatusEmAprovacao Status = "em_aprovacao"
	StatusEmAnalise   Status = "em_analise"
	StatusFinalizada  Status = "finalizada"
	StatusRecusada    Status = "recusada"
)

func (s Status) Terminal() bool {
	return s == StatusFinalizada || s == StatusRecusada
}

// HistoryAction identifica o tipo de evento registrado no histórico da solicitação.
type HistoryAction string

const (
	ActionCriada         HistoryAction = "criada"
	ActionEnviada        HistoryAction = "enviada"
	ActionAprovadaEtapa1 HistoryAction = "aprovada_etapa1"
	ActionAprovadaEtapa2 HistoryAction = "aprovada_etapa2"
	ActionRecusada       HistoryAction = "recusada"
	ActionFinalizada     HistoryAction = "finalizada"
)
