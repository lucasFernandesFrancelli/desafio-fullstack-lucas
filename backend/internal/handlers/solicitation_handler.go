package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/httpx"
	"ekaizen-backend/internal/middleware"
	"ekaizen-backend/internal/models"
	"ekaizen-backend/internal/services"
)

type SolicitationHandler struct {
	service *services.SolicitationService
}

func NewSolicitationHandler(service *services.SolicitationService) *SolicitationHandler {
	return &SolicitationHandler{service: service}
}

func (h *SolicitationHandler) Create(w http.ResponseWriter, r *http.Request) {
	actor, ok := middleware.ActorFrom(r.Context())
	if !ok {
		httpx.Error(w, apperrors.Unauthorized("não autenticado"))
		return
	}

	var input models.DraftInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.Error(w, apperrors.Validation("corpo da requisição inválido", nil))
		return
	}

	sol, err := h.service.CreateDraft(r.Context(), actor, input)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, sol)
}

func (h *SolicitationHandler) List(w http.ResponseWriter, r *http.Request) {
	actor, ok := middleware.ActorFrom(r.Context())
	if !ok {
		httpx.Error(w, apperrors.Unauthorized("não autenticado"))
		return
	}

	filter := models.SolicitationFilter{Query: r.URL.Query().Get("q")}
	if statusParam := r.URL.Query().Get("status"); statusParam != "" {
		status := models.Status(statusParam)
		filter.Status = &status
	}
	if categoryParam := r.URL.Query().Get("categoryId"); categoryParam != "" {
		categoryID, err := uuid.Parse(categoryParam)
		if err != nil {
			httpx.Error(w, apperrors.Validation("categoryId inválido", map[string]string{"categoryId": "deve ser um uuid válido"}))
			return
		}
		filter.CategoryID = &categoryID
	}

	items, err := h.service.List(r.Context(), actor, filter)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	if items == nil {
		items = []models.SolicitationSummary{}
	}
	httpx.JSON(w, http.StatusOK, items)
}

func (h *SolicitationHandler) Get(w http.ResponseWriter, r *http.Request) {
	actor, ok := middleware.ActorFrom(r.Context())
	if !ok {
		httpx.Error(w, apperrors.Unauthorized("não autenticado"))
		return
	}
	id, err := parseID(r)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	detail, err := h.service.GetDetail(r.Context(), actor, id)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, detail)
}

func (h *SolicitationHandler) PatchDraft(w http.ResponseWriter, r *http.Request) {
	actor, ok := middleware.ActorFrom(r.Context())
	if !ok {
		httpx.Error(w, apperrors.Unauthorized("não autenticado"))
		return
	}
	id, err := parseID(r)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	var input models.DraftInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.Error(w, apperrors.Validation("corpo da requisição inválido", nil))
		return
	}

	sol, err := h.service.UpdateDraft(r.Context(), actor, id, input)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, sol)
}

func (h *SolicitationHandler) Submit(w http.ResponseWriter, r *http.Request) {
	actor, ok := middleware.ActorFrom(r.Context())
	if !ok {
		httpx.Error(w, apperrors.Unauthorized("não autenticado"))
		return
	}
	id, err := parseID(r)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	sol, err := h.service.Submit(r.Context(), actor, id)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, sol)
}

func (h *SolicitationHandler) Approve(w http.ResponseWriter, r *http.Request) {
	actor, ok := middleware.ActorFrom(r.Context())
	if !ok {
		httpx.Error(w, apperrors.Unauthorized("não autenticado"))
		return
	}
	id, err := parseID(r)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	sol, err := h.service.Approve(r.Context(), actor, id)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, sol)
}

type rejectRequest struct {
	Reason string `json:"reason"`
}

func (h *SolicitationHandler) Reject(w http.ResponseWriter, r *http.Request) {
	actor, ok := middleware.ActorFrom(r.Context())
	if !ok {
		httpx.Error(w, apperrors.Unauthorized("não autenticado"))
		return
	}
	id, err := parseID(r)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	var req rejectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, apperrors.Validation("corpo da requisição inválido", nil))
		return
	}

	sol, err := h.service.Reject(r.Context(), actor, id, req.Reason)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, sol)
}

func (h *SolicitationHandler) PatchAnalysis(w http.ResponseWriter, r *http.Request) {
	actor, ok := middleware.ActorFrom(r.Context())
	if !ok {
		httpx.Error(w, apperrors.Unauthorized("não autenticado"))
		return
	}
	id, err := parseID(r)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	var input models.AnalysisInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.Error(w, apperrors.Validation("corpo da requisição inválido", nil))
		return
	}

	sol, err := h.service.SaveAnalysis(r.Context(), actor, id, input)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, sol)
}

func (h *SolicitationHandler) Finalize(w http.ResponseWriter, r *http.Request) {
	actor, ok := middleware.ActorFrom(r.Context())
	if !ok {
		httpx.Error(w, apperrors.Unauthorized("não autenticado"))
		return
	}
	id, err := parseID(r)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	var input models.AnalysisInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.Error(w, apperrors.Validation("corpo da requisição inválido", nil))
		return
	}

	sol, err := h.service.Finalize(r.Context(), actor, id, input)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, sol)
}

func (h *SolicitationHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	actor, ok := middleware.ActorFrom(r.Context())
	if !ok {
		httpx.Error(w, apperrors.Unauthorized("não autenticado"))
		return
	}
	id, err := parseID(r)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	history, err := h.service.GetHistory(r.Context(), actor, id)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, history)
}

func parseID(r *http.Request) (uuid.UUID, error) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return uuid.UUID{}, apperrors.Validation("id inválido", map[string]string{"id": "deve ser um uuid válido"})
	}
	return id, nil
}
