package handlers

import (
	"net/http"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/httpx"
	"ekaizen-backend/internal/middleware"
	"ekaizen-backend/internal/services"
)

type DashboardHandler struct {
	service *services.DashboardService
}

func NewDashboardHandler(service *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) Get(w http.ResponseWriter, r *http.Request) {
	actor, ok := middleware.ActorFrom(r.Context())
	if !ok {
		httpx.Error(w, apperrors.Unauthorized("não autenticado"))
		return
	}
	summary, err := h.service.GetSummary(r.Context(), actor)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, summary)
}
