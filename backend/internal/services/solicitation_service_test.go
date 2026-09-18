package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/models"
	"ekaizen-backend/internal/services"
	"ekaizen-backend/internal/testutil"
)

type harness struct {
	sols *testutil.SolicitationRepo
	cats *testutil.CategoryRepo
	hist *testutil.HistoryRepo
	svc  *services.SolicitationService
}

func newHarness() *harness {
	h := &harness{
		sols: &testutil.SolicitationRepo{},
		cats: &testutil.CategoryRepo{},
		hist: &testutil.HistoryRepo{},
	}
	h.svc = services.NewSolicitationService(&testutil.Transactor{}, h.sols, h.cats, h.hist)
	return h
}

func requireAppError(t *testing.T, err error, code apperrors.Code) {
	t.Helper()
	appErr, ok := apperrors.As(err)
	require.True(t, ok, "esperava um *apperrors.AppError, recebeu %v", err)
	assert.Equal(t, code, appErr.Code)
}

func TestCreateDraft_Success(t *testing.T) {
	h := newHarness()
	actor := models.User{ID: uuid.New(), Name: "Ana", Role: models.RoleColaborador}
	title := "Sinalização insuficiente"

	h.sols.On("Create", mock.Anything, mock.MatchedBy(func(s *models.Solicitation) bool {
		return s.Title == title && s.RequesterID == actor.ID && s.Status == models.StatusRascunho
	})).Return(nil)
	h.hist.On("Create", mock.Anything, mock.MatchedBy(func(e *models.HistoryEntry) bool {
		return e.Action == models.ActionCriada && e.ToStatus == models.StatusRascunho
	})).Return(nil)

	sol, err := h.svc.CreateDraft(context.Background(), actor, models.DraftInput{Title: &title})

	require.NoError(t, err)
	assert.Equal(t, title, sol.Title)
	assert.Equal(t, models.StatusRascunho, sol.Status)
}

func TestCreateDraft_MissingTitle(t *testing.T) {
	h := newHarness()
	actor := models.User{ID: uuid.New()}

	_, err := h.svc.CreateDraft(context.Background(), actor, models.DraftInput{})

	requireAppError(t, err, apperrors.CodeValidation)
}

func TestUpdateDraft_Success(t *testing.T) {
	h := newHarness()
	actor := models.User{ID: uuid.New()}
	solID := uuid.New()
	newTitle := "Novo título"
	existing := &models.Solicitation{ID: solID, RequesterID: actor.ID, Status: models.StatusRascunho, Title: "Antigo"}

	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)
	h.sols.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Solicitation) bool {
		return s.Title == newTitle
	})).Return(nil)
	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)

	sol, err := h.svc.UpdateDraft(context.Background(), actor, solID, models.DraftInput{Title: &newTitle})

	require.NoError(t, err)
	assert.Equal(t, newTitle, sol.Title)
}

func TestUpdateDraft_Forbidden_NotOwner(t *testing.T) {
	h := newHarness()
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, RequesterID: uuid.New(), Status: models.StatusRascunho}
	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.UpdateDraft(context.Background(), models.User{ID: uuid.New()}, solID, models.DraftInput{})

	requireAppError(t, err, apperrors.CodeForbidden)
}

func TestUpdateDraft_Conflict_NotDraft(t *testing.T) {
	h := newHarness()
	actor := models.User{ID: uuid.New()}
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, RequesterID: actor.ID, Status: models.StatusEmAprovacao}
	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.UpdateDraft(context.Background(), actor, solID, models.DraftInput{})

	requireAppError(t, err, apperrors.CodeConflict)
}

func TestSubmit_Success(t *testing.T) {
	h := newHarness()
	actor := models.User{ID: uuid.New()}
	catID := uuid.New()
	solID := uuid.New()
	existing := &models.Solicitation{
		ID: solID, RequesterID: actor.ID, Status: models.StatusRascunho,
		Title: "T", ProblemDescription: "P", ProposedImprovement: "M", CategoryID: &catID, Location: "L",
	}

	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)
	h.sols.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Solicitation) bool {
		return s.Status == models.StatusEmAprovacao && s.CurrentApprovalStep != nil && *s.CurrentApprovalStep == 1
	})).Return(nil)
	h.hist.On("Create", mock.Anything, mock.MatchedBy(func(e *models.HistoryEntry) bool {
		return e.Action == models.ActionEnviada
	})).Return(nil)
	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.Submit(context.Background(), actor, solID)

	require.NoError(t, err)
}

func TestSubmit_ValidationMissingFields(t *testing.T) {
	h := newHarness()
	actor := models.User{ID: uuid.New()}
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, RequesterID: actor.ID, Status: models.StatusRascunho, Title: "Só título"}
	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.Submit(context.Background(), actor, solID)

	requireAppError(t, err, apperrors.CodeValidation)
	appErr, _ := apperrors.As(err)
	assert.Contains(t, appErr.Fields, "categoryId")
}

func TestSubmit_Forbidden_NotOwner(t *testing.T) {
	h := newHarness()
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, RequesterID: uuid.New(), Status: models.StatusRascunho}
	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.Submit(context.Background(), models.User{ID: uuid.New()}, solID)

	requireAppError(t, err, apperrors.CodeForbidden)
}

func TestSubmit_Conflict_NotDraft(t *testing.T) {
	h := newHarness()
	actor := models.User{ID: uuid.New()}
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, RequesterID: actor.ID, Status: models.StatusFinalizada}
	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.Submit(context.Background(), actor, solID)

	requireAppError(t, err, apperrors.CodeConflict)
}

func TestApprove_Step1_Success(t *testing.T) {
	h := newHarness()
	approver := models.User{ID: uuid.New()}
	catID := uuid.New()
	solID := uuid.New()
	step1 := int16(1)
	existing := &models.Solicitation{ID: solID, Status: models.StatusEmAprovacao, CategoryID: &catID, CurrentApprovalStep: &step1}

	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)
	h.cats.On("GetApproverOrder", mock.Anything, catID, approver.ID).Return(&step1, nil)
	h.sols.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Solicitation) bool {
		return s.Status == models.StatusEmAprovacao && s.CurrentApprovalStep != nil && *s.CurrentApprovalStep == 2
	})).Return(nil)
	h.hist.On("Create", mock.Anything, mock.MatchedBy(func(e *models.HistoryEntry) bool {
		return e.Action == models.ActionAprovadaEtapa1
	})).Return(nil)
	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.Approve(context.Background(), approver, solID)

	require.NoError(t, err)
}

func TestApprove_Step2_Success_MovesToAnalysis(t *testing.T) {
	h := newHarness()
	approver := models.User{ID: uuid.New()}
	catID := uuid.New()
	solID := uuid.New()
	step2 := int16(2)
	existing := &models.Solicitation{ID: solID, Status: models.StatusEmAprovacao, CategoryID: &catID, CurrentApprovalStep: &step2}

	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)
	h.cats.On("GetApproverOrder", mock.Anything, catID, approver.ID).Return(&step2, nil)
	h.sols.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Solicitation) bool {
		return s.Status == models.StatusEmAnalise && s.CurrentApprovalStep == nil
	})).Return(nil)
	h.hist.On("Create", mock.Anything, mock.MatchedBy(func(e *models.HistoryEntry) bool {
		return e.Action == models.ActionAprovadaEtapa2
	})).Return(nil)
	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.Approve(context.Background(), approver, solID)

	require.NoError(t, err)
}

func TestApprove_Forbidden_WrongTurn(t *testing.T) {
	h := newHarness()
	actor := models.User{ID: uuid.New()}
	catID := uuid.New()
	solID := uuid.New()
	step1 := int16(1)
	existing := &models.Solicitation{ID: solID, Status: models.StatusEmAprovacao, CategoryID: &catID, CurrentApprovalStep: &step1}

	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)
	h.cats.On("GetApproverOrder", mock.Anything, catID, actor.ID).Return(nil, nil)

	_, err := h.svc.Approve(context.Background(), actor, solID)

	requireAppError(t, err, apperrors.CodeForbidden)
}

func TestApprove_Conflict_WrongStatus(t *testing.T) {
	h := newHarness()
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, Status: models.StatusFinalizada}
	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.Approve(context.Background(), models.User{ID: uuid.New()}, solID)

	requireAppError(t, err, apperrors.CodeConflict)
}

func TestReject_Success(t *testing.T) {
	h := newHarness()
	actor := models.User{ID: uuid.New()}
	catID := uuid.New()
	solID := uuid.New()
	step1 := int16(1)
	existing := &models.Solicitation{ID: solID, Status: models.StatusEmAprovacao, CategoryID: &catID, CurrentApprovalStep: &step1}

	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)
	h.cats.On("GetApproverOrder", mock.Anything, catID, actor.ID).Return(&step1, nil)
	h.sols.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Solicitation) bool {
		return s.Status == models.StatusRecusada && s.RejectionReason == "Motivo válido"
	})).Return(nil)
	h.hist.On("Create", mock.Anything, mock.MatchedBy(func(e *models.HistoryEntry) bool {
		return e.Action == models.ActionRecusada && e.Comment == "Motivo válido"
	})).Return(nil)
	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.Reject(context.Background(), actor, solID, "Motivo válido")

	require.NoError(t, err)
}

func TestReject_ValidationEmptyReason(t *testing.T) {
	h := newHarness()

	_, err := h.svc.Reject(context.Background(), models.User{ID: uuid.New()}, uuid.New(), "   ")

	requireAppError(t, err, apperrors.CodeValidation)
}

func TestReject_Forbidden_WrongApprover(t *testing.T) {
	h := newHarness()
	actor := models.User{ID: uuid.New()}
	catID := uuid.New()
	solID := uuid.New()
	step1, step2 := int16(1), int16(2)
	existing := &models.Solicitation{ID: solID, Status: models.StatusEmAprovacao, CategoryID: &catID, CurrentApprovalStep: &step2}

	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)
	h.cats.On("GetApproverOrder", mock.Anything, catID, actor.ID).Return(&step1, nil)

	_, err := h.svc.Reject(context.Background(), actor, solID, "motivo")

	requireAppError(t, err, apperrors.CodeForbidden)
}

func TestSaveAnalysis_Success(t *testing.T) {
	h := newHarness()
	analyst := models.User{ID: uuid.New(), Role: models.RoleAnalista}
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, Status: models.StatusEmAnalise}
	sev := int16(3)

	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)
	h.sols.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Solicitation) bool {
		return s.Severity != nil && *s.Severity == 3
	})).Return(nil)
	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.SaveAnalysis(context.Background(), analyst, solID, models.AnalysisInput{Severity: &sev})

	require.NoError(t, err)
}

func TestSaveAnalysis_Forbidden_NotAnalyst(t *testing.T) {
	h := newHarness()

	_, err := h.svc.SaveAnalysis(context.Background(), models.User{ID: uuid.New(), Role: models.RoleColaborador}, uuid.New(), models.AnalysisInput{})

	requireAppError(t, err, apperrors.CodeForbidden)
}

func TestSaveAnalysis_Conflict_WrongStatus(t *testing.T) {
	h := newHarness()
	analyst := models.User{ID: uuid.New(), Role: models.RoleAnalista}
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, Status: models.StatusEmAprovacao}
	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.SaveAnalysis(context.Background(), analyst, solID, models.AnalysisInput{})

	requireAppError(t, err, apperrors.CodeConflict)
}

func TestFinalize_Success(t *testing.T) {
	h := newHarness()
	analyst := models.User{ID: uuid.New(), Role: models.RoleAnalista}
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, Status: models.StatusEmAnalise}
	sev, urg, trend := int16(4), int16(3), int16(5)
	notes := "Parecer completo"

	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)
	h.sols.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Solicitation) bool {
		return s.Status == models.StatusFinalizada
	})).Return(nil)
	h.hist.On("Create", mock.Anything, mock.MatchedBy(func(e *models.HistoryEntry) bool {
		return e.Action == models.ActionFinalizada
	})).Return(nil)
	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.Finalize(context.Background(), analyst, solID, models.AnalysisInput{
		Severity: &sev, Urgency: &urg, Trend: &trend, AnalysisNotes: &notes,
	})

	require.NoError(t, err)
}

func TestFinalize_ValidationMissingFields(t *testing.T) {
	h := newHarness()
	analyst := models.User{ID: uuid.New(), Role: models.RoleAnalista}
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, Status: models.StatusEmAnalise}
	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.Finalize(context.Background(), analyst, solID, models.AnalysisInput{})

	requireAppError(t, err, apperrors.CodeValidation)
	appErr, _ := apperrors.As(err)
	assert.Contains(t, appErr.Fields, "severity")
}

func TestFinalize_Forbidden_NotAnalyst(t *testing.T) {
	h := newHarness()

	_, err := h.svc.Finalize(context.Background(), models.User{ID: uuid.New(), Role: models.RoleGestor}, uuid.New(), models.AnalysisInput{})

	requireAppError(t, err, apperrors.CodeForbidden)
}

func TestGetDetail_DraftHiddenFromOthers(t *testing.T) {
	h := newHarness()
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, RequesterID: uuid.New(), Status: models.StatusRascunho}
	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.GetDetail(context.Background(), models.User{ID: uuid.New(), Role: models.RoleColaborador}, solID)

	requireAppError(t, err, apperrors.CodeNotFound)
}

func TestGetDetail_DraftVisibleToOwnerAndManager(t *testing.T) {
	h := newHarness()
	ownerID := uuid.New()
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, RequesterID: ownerID, RequesterName: "Dono", Status: models.StatusRascunho}

	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)
	h.hist.On("ListBySolicitation", mock.Anything, solID).Return([]models.HistoryEntry{}, nil)

	owner := models.User{ID: ownerID, Role: models.RoleColaborador}
	detail, err := h.svc.GetDetail(context.Background(), owner, solID)
	require.NoError(t, err)
	assert.True(t, detail.Permissions.CanEdit)
	assert.True(t, detail.Permissions.CanSubmit)

	manager := models.User{ID: uuid.New(), Role: models.RoleGestor}
	detail2, err := h.svc.GetDetail(context.Background(), manager, solID)
	require.NoError(t, err)
	assert.False(t, detail2.Permissions.CanEdit)
}

func TestGetDetail_PermissionsForApprover(t *testing.T) {
	h := newHarness()
	approver := models.User{ID: uuid.New()}
	catID := uuid.New()
	solID := uuid.New()
	step1 := int16(1)
	existing := &models.Solicitation{ID: solID, Status: models.StatusEmAprovacao, CategoryID: &catID, CurrentApprovalStep: &step1}

	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)
	h.cats.On("GetApproverByOrder", mock.Anything, catID, step1).Return(&models.User{ID: approver.ID, Name: "Aprovador"}, nil)
	h.cats.On("GetApproverOrder", mock.Anything, catID, approver.ID).Return(&step1, nil)
	h.hist.On("ListBySolicitation", mock.Anything, solID).Return([]models.HistoryEntry{}, nil)

	detail, err := h.svc.GetDetail(context.Background(), approver, solID)

	require.NoError(t, err)
	assert.True(t, detail.Permissions.CanApprove)
	assert.True(t, detail.Permissions.CanReject)
	assert.Equal(t, "Aprovador", detail.PendingActor.Name)
}

func TestUpdateDraft_AllFields(t *testing.T) {
	h := newHarness()
	actor := models.User{ID: uuid.New()}
	solID := uuid.New()
	catID := uuid.New()
	existing := &models.Solicitation{ID: solID, RequesterID: actor.ID, Status: models.StatusRascunho}

	title, problem, improvement, location := "T", "P", "M", "L"
	input := models.DraftInput{
		Title: &title, ProblemDescription: &problem, ProposedImprovement: &improvement,
		CategoryID: &catID, Location: &location,
	}

	h.sols.On("GetForUpdate", mock.Anything, solID).Return(existing, nil)
	h.sols.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Solicitation) bool {
		return s.Title == title && s.ProblemDescription == problem && s.ProposedImprovement == improvement &&
			s.CategoryID != nil && *s.CategoryID == catID && s.Location == location
	})).Return(nil)
	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.UpdateDraft(context.Background(), actor, solID, input)
	require.NoError(t, err)
}

func TestGetDetail_TerminalStatus_NoPendingActorNoPermissions(t *testing.T) {
	h := newHarness()
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, RequesterID: uuid.New(), Status: models.StatusFinalizada}

	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)
	h.hist.On("ListBySolicitation", mock.Anything, solID).Return([]models.HistoryEntry{}, nil)

	detail, err := h.svc.GetDetail(context.Background(), models.User{ID: uuid.New(), Role: models.RoleGestor}, solID)

	require.NoError(t, err)
	assert.Equal(t, "none", detail.PendingActor.Kind)
	assert.False(t, detail.Permissions.CanEdit)
	assert.False(t, detail.Permissions.CanApprove)
	assert.False(t, detail.Permissions.CanAnalyze)
}

func TestGetDetail_PropagatesPendingActorError(t *testing.T) {
	h := newHarness()
	catID := uuid.New()
	solID := uuid.New()
	step1 := int16(1)
	existing := &models.Solicitation{ID: solID, Status: models.StatusEmAprovacao, CategoryID: &catID, CurrentApprovalStep: &step1}

	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)
	h.cats.On("GetApproverByOrder", mock.Anything, catID, step1).Return(nil, errors.New("falha no banco"))

	_, err := h.svc.GetDetail(context.Background(), models.User{ID: uuid.New()}, solID)

	require.Error(t, err)
}

func TestGetDetail_PropagatesPermissionsError(t *testing.T) {
	h := newHarness()
	catID := uuid.New()
	solID := uuid.New()
	step1 := int16(1)
	approverID := uuid.New()
	existing := &models.Solicitation{ID: solID, Status: models.StatusEmAprovacao, CategoryID: &catID, CurrentApprovalStep: &step1}

	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)
	h.cats.On("GetApproverByOrder", mock.Anything, catID, step1).Return(&models.User{ID: approverID, Name: "Aprovador"}, nil)
	h.cats.On("GetApproverOrder", mock.Anything, catID, approverID).Return(nil, errors.New("falha no banco"))

	_, err := h.svc.GetDetail(context.Background(), models.User{ID: approverID}, solID)

	require.Error(t, err)
}

func TestGetDetail_PropagatesHistoryError(t *testing.T) {
	h := newHarness()
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, Status: models.StatusFinalizada}

	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)
	h.hist.On("ListBySolicitation", mock.Anything, solID).Return(nil, errors.New("falha no banco"))

	_, err := h.svc.GetDetail(context.Background(), models.User{ID: uuid.New()}, solID)

	require.Error(t, err)
}

func TestGetHistory_PropagatesHistoryError(t *testing.T) {
	h := newHarness()
	solID := uuid.New()
	actor := models.User{ID: uuid.New()}
	existing := &models.Solicitation{ID: solID, RequesterID: actor.ID, Status: models.StatusRascunho}

	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)
	h.hist.On("ListBySolicitation", mock.Anything, solID).Return(nil, errors.New("falha no banco"))

	_, err := h.svc.GetHistory(context.Background(), actor, solID)

	require.Error(t, err)
}

func TestGetHistory_Success(t *testing.T) {
	h := newHarness()
	actor := models.User{ID: uuid.New(), Role: models.RoleColaborador}
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, RequesterID: actor.ID, Status: models.StatusEmAnalise}

	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)
	h.hist.On("ListBySolicitation", mock.Anything, solID).Return([]models.HistoryEntry{
		{ID: uuid.New(), SolicitationID: solID, ToStatus: models.StatusRascunho, Action: models.ActionCriada},
	}, nil)

	entries, err := h.svc.GetHistory(context.Background(), actor, solID)

	require.NoError(t, err)
	require.Len(t, entries, 1)
}

func TestGetHistory_HiddenDraft(t *testing.T) {
	h := newHarness()
	solID := uuid.New()
	existing := &models.Solicitation{ID: solID, RequesterID: uuid.New(), Status: models.StatusRascunho}
	h.sols.On("GetByID", mock.Anything, solID).Return(existing, nil)

	_, err := h.svc.GetHistory(context.Background(), models.User{ID: uuid.New(), Role: models.RoleColaborador}, solID)

	requireAppError(t, err, apperrors.CodeNotFound)
}

func TestList_PropagatesPendingActorError(t *testing.T) {
	h := newHarness()
	catID := uuid.New()
	step1 := int16(1)
	item := models.Solicitation{ID: uuid.New(), Status: models.StatusEmAprovacao, CategoryID: &catID, CurrentApprovalStep: &step1}

	h.sols.On("List", mock.Anything, mock.Anything).Return([]models.Solicitation{item}, nil)
	h.cats.On("GetApproverByOrder", mock.Anything, catID, step1).Return(nil, errors.New("falha no banco"))

	_, err := h.svc.List(context.Background(), models.User{ID: uuid.New(), Role: models.RoleGestor}, models.SolicitationFilter{})

	require.Error(t, err)
}

func TestGetDetail_PropagatesGetByIDError(t *testing.T) {
	h := newHarness()
	solID := uuid.New()

	h.sols.On("GetByID", mock.Anything, solID).Return(nil, errors.New("falha no banco"))

	_, err := h.svc.GetDetail(context.Background(), models.User{ID: uuid.New()}, solID)

	require.Error(t, err)
}

func TestGetHistory_PropagatesGetByIDError(t *testing.T) {
	h := newHarness()
	solID := uuid.New()

	h.sols.On("GetByID", mock.Anything, solID).Return(nil, errors.New("falha no banco"))

	_, err := h.svc.GetHistory(context.Background(), models.User{ID: uuid.New()}, solID)

	require.Error(t, err)
}

func TestList_PropagatesRepositoryError(t *testing.T) {
	h := newHarness()

	h.sols.On("List", mock.Anything, mock.Anything).Return(nil, errors.New("falha no banco"))

	_, err := h.svc.List(context.Background(), models.User{ID: uuid.New(), Role: models.RoleGestor}, models.SolicitationFilter{})

	require.Error(t, err)
}

func TestList_FiltersOutDraftsNotOwned(t *testing.T) {
	h := newHarness()
	owner := uuid.New()
	viewer := models.User{ID: uuid.New(), Role: models.RoleColaborador}

	draftOfSomeoneElse := models.Solicitation{ID: uuid.New(), RequesterID: owner, Status: models.StatusRascunho}
	visible := models.Solicitation{ID: uuid.New(), RequesterID: owner, Status: models.StatusEmAnalise}

	h.sols.On("List", mock.Anything, mock.Anything).Return([]models.Solicitation{draftOfSomeoneElse, visible}, nil)

	result, err := h.svc.List(context.Background(), viewer, models.SolicitationFilter{})

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, visible.ID, result[0].ID)
}
