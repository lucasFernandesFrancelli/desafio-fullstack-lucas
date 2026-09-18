// Package services contém a lógica de negócio. SolicitationService é a
// ÚNICA porta de entrada para criar ou mudar o estado de uma solicitação —
// handlers HTTP nunca tocam o repositório diretamente, o que garante que
// Lista, Kanban e Detalhe sempre respeitem exatamente as mesmas regras.
package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/models"
	"ekaizen-backend/internal/repository"
)

type SolicitationService struct {
	tx            repository.Transactor
	solicitations repository.SolicitationRepository
	categories    repository.CategoryRepository
	history       repository.HistoryRepository
}

func NewSolicitationService(
	tx repository.Transactor,
	solicitations repository.SolicitationRepository,
	categories repository.CategoryRepository,
	history repository.HistoryRepository,
) *SolicitationService {
	return &SolicitationService{tx: tx, solicitations: solicitations, categories: categories, history: history}
}

func (s *SolicitationService) CreateDraft(ctx context.Context, actor models.User, input models.DraftInput) (*models.Solicitation, error) {
	sol := &models.Solicitation{
		ID:          uuid.New(),
		RequesterID: actor.ID,
		Status:      models.StatusRascunho,
	}
	applyDraftInput(sol, input)
	if strings.TrimSpace(sol.Title) == "" {
		return nil, apperrors.Validation("título é obrigatório", map[string]string{"title": "obrigatório"})
	}

	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		if err := s.solicitations.Create(ctx, sol); err != nil {
			return err
		}
		return s.history.Create(ctx, &models.HistoryEntry{
			ID: uuid.New(), SolicitationID: sol.ID,
			ToStatus: models.StatusRascunho, Action: models.ActionCriada, ActorID: actor.ID,
		})
	})
	if err != nil {
		return nil, err
	}
	sol.RequesterName = actor.Name
	return sol, nil
}

func (s *SolicitationService) UpdateDraft(ctx context.Context, actor models.User, id uuid.UUID, input models.DraftInput) (*models.Solicitation, error) {
	var result *models.Solicitation
	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		sol, err := s.solicitations.GetForUpdate(ctx, id)
		if err != nil {
			return err
		}
		if sol.RequesterID != actor.ID {
			return apperrors.Forbidden("apenas o solicitante pode editar este rascunho")
		}
		if sol.Status != models.StatusRascunho {
			return apperrors.Conflict("solicitação não está mais em rascunho")
		}
		applyDraftInput(sol, input)
		if err := s.solicitations.Update(ctx, sol); err != nil {
			return err
		}
		result = sol
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.solicitations.GetByID(ctx, result.ID)
}

func (s *SolicitationService) Submit(ctx context.Context, actor models.User, id uuid.UUID) (*models.Solicitation, error) {
	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		sol, err := s.solicitations.GetForUpdate(ctx, id)
		if err != nil {
			return err
		}
		if sol.RequesterID != actor.ID {
			return apperrors.Forbidden("apenas o solicitante pode enviar esta solicitação")
		}
		if sol.Status != models.StatusRascunho {
			return apperrors.Conflict("solicitação não está em rascunho")
		}

		fields := map[string]string{}
		if strings.TrimSpace(sol.Title) == "" {
			fields["title"] = "obrigatório"
		}
		if strings.TrimSpace(sol.ProblemDescription) == "" {
			fields["problemDescription"] = "obrigatório"
		}
		if strings.TrimSpace(sol.ProposedImprovement) == "" {
			fields["proposedImprovement"] = "obrigatório"
		}
		if sol.CategoryID == nil {
			fields["categoryId"] = "obrigatório"
		}
		if strings.TrimSpace(sol.Location) == "" {
			fields["location"] = "obrigatório"
		}
		if len(fields) > 0 {
			return apperrors.Validation("preencha todos os campos antes de enviar", fields)
		}

		fromStatus := sol.Status
		step := int16(1)
		sol.CurrentApprovalStep = &step
		sol.Status = models.StatusEmAprovacao
		sol.LastTransitionAt = time.Now()
		if err := s.solicitations.Update(ctx, sol); err != nil {
			return err
		}
		return s.history.Create(ctx, &models.HistoryEntry{
			ID: uuid.New(), SolicitationID: sol.ID,
			FromStatus: statusPtr(fromStatus), ToStatus: sol.Status,
			Action: models.ActionEnviada, ActorID: actor.ID,
		})
	})
	if err != nil {
		return nil, err
	}
	return s.solicitations.GetByID(ctx, id)
}

func (s *SolicitationService) Approve(ctx context.Context, actor models.User, id uuid.UUID) (*models.Solicitation, error) {
	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		sol, err := s.solicitations.GetForUpdate(ctx, id)
		if err != nil {
			return err
		}
		if sol.Status != models.StatusEmAprovacao {
			return apperrors.Conflict("solicitação não está em aprovação")
		}
		if sol.CategoryID == nil || sol.CurrentApprovalStep == nil {
			return apperrors.Conflict("solicitação sem categoria ou etapa de aprovação definida")
		}

		order, err := s.categories.GetApproverOrder(ctx, *sol.CategoryID, actor.ID)
		if err != nil {
			return err
		}
		if order == nil || *order != *sol.CurrentApprovalStep {
			return apperrors.Forbidden("não é a sua vez de aprovar esta solicitação")
		}

		fromStatus := sol.Status
		currentStep := *sol.CurrentApprovalStep
		sol.LastTransitionAt = time.Now()

		if currentStep == 1 {
			nextStep := int16(2)
			sol.CurrentApprovalStep = &nextStep
			if err := s.solicitations.Update(ctx, sol); err != nil {
				return err
			}
			return s.history.Create(ctx, &models.HistoryEntry{
				ID: uuid.New(), SolicitationID: sol.ID,
				FromStatus: statusPtr(fromStatus), ToStatus: sol.Status,
				Action: models.ActionAprovadaEtapa1, ActorID: actor.ID, ApprovalStep: int16Ptr(1),
			})
		}

		sol.Status = models.StatusEmAnalise
		sol.CurrentApprovalStep = nil
		if err := s.solicitations.Update(ctx, sol); err != nil {
			return err
		}
		return s.history.Create(ctx, &models.HistoryEntry{
			ID: uuid.New(), SolicitationID: sol.ID,
			FromStatus: statusPtr(fromStatus), ToStatus: sol.Status,
			Action: models.ActionAprovadaEtapa2, ActorID: actor.ID, ApprovalStep: int16Ptr(2),
		})
	})
	if err != nil {
		return nil, err
	}
	return s.solicitations.GetByID(ctx, id)
}

func (s *SolicitationService) Reject(ctx context.Context, actor models.User, id uuid.UUID, reason string) (*models.Solicitation, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, apperrors.Validation("motivo da recusa é obrigatório", map[string]string{"reason": "obrigatório"})
	}

	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		sol, err := s.solicitations.GetForUpdate(ctx, id)
		if err != nil {
			return err
		}
		if sol.Status != models.StatusEmAprovacao {
			return apperrors.Conflict("solicitação não está em aprovação")
		}
		if sol.CategoryID == nil || sol.CurrentApprovalStep == nil {
			return apperrors.Conflict("solicitação sem categoria ou etapa de aprovação definida")
		}

		order, err := s.categories.GetApproverOrder(ctx, *sol.CategoryID, actor.ID)
		if err != nil {
			return err
		}
		if order == nil || *order != *sol.CurrentApprovalStep {
			return apperrors.Forbidden("não é a sua vez de decidir esta solicitação")
		}

		fromStatus := sol.Status
		step := *sol.CurrentApprovalStep
		sol.Status = models.StatusRecusada
		sol.CurrentApprovalStep = nil
		sol.RejectionReason = reason
		sol.LastTransitionAt = time.Now()
		if err := s.solicitations.Update(ctx, sol); err != nil {
			return err
		}
		return s.history.Create(ctx, &models.HistoryEntry{
			ID: uuid.New(), SolicitationID: sol.ID,
			FromStatus: statusPtr(fromStatus), ToStatus: sol.Status,
			Action: models.ActionRecusada, ActorID: actor.ID, ApprovalStep: &step, Comment: reason,
		})
	})
	if err != nil {
		return nil, err
	}
	return s.solicitations.GetByID(ctx, id)
}

func (s *SolicitationService) SaveAnalysis(ctx context.Context, actor models.User, id uuid.UUID, input models.AnalysisInput) (*models.Solicitation, error) {
	if actor.Role != models.RoleAnalista {
		return nil, apperrors.Forbidden("apenas analistas podem registrar parecer")
	}

	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		sol, err := s.solicitations.GetForUpdate(ctx, id)
		if err != nil {
			return err
		}
		if sol.Status != models.StatusEmAnalise {
			return apperrors.Conflict("solicitação não está em análise")
		}
		if err := applyAnalysisInput(sol, input); err != nil {
			return err
		}
		return s.solicitations.Update(ctx, sol)
	})
	if err != nil {
		return nil, err
	}
	return s.solicitations.GetByID(ctx, id)
}

func (s *SolicitationService) Finalize(ctx context.Context, actor models.User, id uuid.UUID, input models.AnalysisInput) (*models.Solicitation, error) {
	if actor.Role != models.RoleAnalista {
		return nil, apperrors.Forbidden("apenas analistas podem finalizar")
	}

	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		sol, err := s.solicitations.GetForUpdate(ctx, id)
		if err != nil {
			return err
		}
		if sol.Status != models.StatusEmAnalise {
			return apperrors.Conflict("solicitação não está em análise")
		}
		if err := applyAnalysisInput(sol, input); err != nil {
			return err
		}

		fields := map[string]string{}
		if sol.Severity == nil {
			fields["severity"] = "obrigatório"
		}
		if sol.Urgency == nil {
			fields["urgency"] = "obrigatório"
		}
		if sol.Trend == nil {
			fields["trend"] = "obrigatório"
		}
		if strings.TrimSpace(sol.AnalysisNotes) == "" {
			fields["analysisNotes"] = "obrigatório"
		}
		if len(fields) > 0 {
			return apperrors.Validation("preencha o parecer e as três notas antes de finalizar", fields)
		}

		fromStatus := sol.Status
		sol.Status = models.StatusFinalizada
		sol.LastTransitionAt = time.Now()
		if err := s.solicitations.Update(ctx, sol); err != nil {
			return err
		}
		return s.history.Create(ctx, &models.HistoryEntry{
			ID: uuid.New(), SolicitationID: sol.ID,
			FromStatus: statusPtr(fromStatus), ToStatus: sol.Status,
			Action: models.ActionFinalizada, ActorID: actor.ID, Comment: sol.AnalysisNotes,
			Severity: sol.Severity, Urgency: sol.Urgency, Trend: sol.Trend,
		})
	})
	if err != nil {
		return nil, err
	}
	return s.solicitations.GetByID(ctx, id)
}

func (s *SolicitationService) GetDetail(ctx context.Context, actor models.User, id uuid.UUID) (*models.SolicitationDetail, error) {
	sol, err := s.solicitations.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !canView(actor, sol) {
		return nil, apperrors.NotFound("solicitação não encontrada")
	}

	pendingActor, err := s.pendingActorFor(ctx, sol)
	if err != nil {
		return nil, err
	}
	perms, err := s.permissionsFor(ctx, actor, sol)
	if err != nil {
		return nil, err
	}
	history, err := s.history.ListBySolicitation(ctx, id)
	if err != nil {
		return nil, err
	}
	if history == nil {
		history = []models.HistoryEntry{}
	}

	return &models.SolicitationDetail{
		Solicitation: *sol,
		PendingActor: pendingActor,
		Permissions:  perms,
		History:      history,
	}, nil
}

func (s *SolicitationService) GetHistory(ctx context.Context, actor models.User, id uuid.UUID) ([]models.HistoryEntry, error) {
	sol, err := s.solicitations.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !canView(actor, sol) {
		return nil, apperrors.NotFound("solicitação não encontrada")
	}
	history, err := s.history.ListBySolicitation(ctx, id)
	if err != nil {
		return nil, err
	}
	if history == nil {
		history = []models.HistoryEntry{}
	}
	return history, nil
}

func (s *SolicitationService) List(ctx context.Context, actor models.User, filter models.SolicitationFilter) ([]models.SolicitationSummary, error) {
	items, err := s.solicitations.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := make([]models.SolicitationSummary, 0, len(items))
	for i := range items {
		sol := items[i]
		if !canView(actor, &sol) {
			continue
		}
		pendingActor, err := s.pendingActorFor(ctx, &sol)
		if err != nil {
			return nil, err
		}
		result = append(result, models.SolicitationSummary{
			ID: sol.ID, Title: sol.Title, CategoryName: sol.CategoryName, Status: sol.Status,
			CurrentApprovalStep: sol.CurrentApprovalStep, RequesterName: sol.RequesterName,
			Priority: sol.Priority, PendingActorName: pendingActor.Name,
			CreatedAt: sol.CreatedAt, UpdatedAt: sol.UpdatedAt,
		})
	}
	return result, nil
}

// canView aplica a regra de visibilidade: rascunhos só são visíveis para o
// dono e para o gestor; os demais estados são visíveis a todos os perfis.
func canView(actor models.User, sol *models.Solicitation) bool {
	if sol.Status != models.StatusRascunho {
		return true
	}
	return actor.Role == models.RoleGestor || actor.ID == sol.RequesterID
}

func (s *SolicitationService) pendingActorFor(ctx context.Context, sol *models.Solicitation) (models.PendingActor, error) {
	switch sol.Status {
	case models.StatusRascunho:
		id := sol.RequesterID.String()
		return models.PendingActor{Kind: "user", ID: &id, Name: sol.RequesterName}, nil
	case models.StatusEmAprovacao:
		if sol.CategoryID == nil || sol.CurrentApprovalStep == nil {
			return models.PendingActor{Kind: "none"}, nil
		}
		approver, err := s.categories.GetApproverByOrder(ctx, *sol.CategoryID, *sol.CurrentApprovalStep)
		if err != nil {
			return models.PendingActor{}, err
		}
		id := approver.ID.String()
		return models.PendingActor{Kind: "user", ID: &id, Name: approver.Name}, nil
	case models.StatusEmAnalise:
		return models.PendingActor{Kind: "role", Name: "Qualquer analista"}, nil
	default:
		return models.PendingActor{Kind: "none"}, nil
	}
}

func (s *SolicitationService) permissionsFor(ctx context.Context, actor models.User, sol *models.Solicitation) (models.Permissions, error) {
	perms := models.Permissions{}
	switch sol.Status {
	case models.StatusRascunho:
		if actor.ID == sol.RequesterID {
			perms.CanEdit = true
			perms.CanSubmit = true
		}
	case models.StatusEmAprovacao:
		if sol.CategoryID != nil && sol.CurrentApprovalStep != nil {
			order, err := s.categories.GetApproverOrder(ctx, *sol.CategoryID, actor.ID)
			if err != nil {
				return perms, err
			}
			if order != nil && *order == *sol.CurrentApprovalStep {
				perms.CanApprove = true
				perms.CanReject = true
			}
		}
	case models.StatusEmAnalise:
		if actor.Role == models.RoleAnalista {
			perms.CanAnalyze = true
			perms.CanFinalize = true
		}
	}
	return perms, nil
}

func applyDraftInput(sol *models.Solicitation, input models.DraftInput) {
	if input.Title != nil {
		sol.Title = strings.TrimSpace(*input.Title)
	}
	if input.ProblemDescription != nil {
		sol.ProblemDescription = strings.TrimSpace(*input.ProblemDescription)
	}
	if input.ProposedImprovement != nil {
		sol.ProposedImprovement = strings.TrimSpace(*input.ProposedImprovement)
	}
	if input.CategoryID != nil {
		sol.CategoryID = input.CategoryID
	}
	if input.Location != nil {
		sol.Location = strings.TrimSpace(*input.Location)
	}
}

func applyAnalysisInput(sol *models.Solicitation, input models.AnalysisInput) error {
	if input.Severity != nil {
		if *input.Severity < 1 || *input.Severity > 5 {
			return apperrors.Validation("gravidade deve estar entre 1 e 5", map[string]string{"severity": "deve estar entre 1 e 5"})
		}
		sol.Severity = input.Severity
	}
	if input.Urgency != nil {
		if *input.Urgency < 1 || *input.Urgency > 5 {
			return apperrors.Validation("urgência deve estar entre 1 e 5", map[string]string{"urgency": "deve estar entre 1 e 5"})
		}
		sol.Urgency = input.Urgency
	}
	if input.Trend != nil {
		if *input.Trend < 1 || *input.Trend > 5 {
			return apperrors.Validation("tendência deve estar entre 1 e 5", map[string]string{"trend": "deve estar entre 1 e 5"})
		}
		sol.Trend = input.Trend
	}
	if input.AnalysisNotes != nil {
		sol.AnalysisNotes = strings.TrimSpace(*input.AnalysisNotes)
	}
	return nil
}

func statusPtr(s models.Status) *models.Status { return &s }
func int16Ptr(v int16) *int16                  { return &v }
