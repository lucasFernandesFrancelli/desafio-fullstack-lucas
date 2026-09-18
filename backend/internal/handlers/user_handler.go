package handlers

import (
	"encoding/json"
	"net/http"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/httpx"
	"ekaizen-backend/internal/middleware"
	"ekaizen-backend/internal/models"
	"ekaizen-backend/internal/services"
)

type UserHandler struct {
	admin *services.AdminService
}

func NewUserHandler(admin *services.AdminService) *UserHandler {
	return &UserHandler{admin: admin}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	actor, ok := middleware.ActorFrom(r.Context())
	if !ok {
		httpx.Error(w, apperrors.Unauthorized("não autenticado"))
		return
	}

	users, err := h.admin.ListUsers(r.Context(), actor)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, users)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	actor, ok := middleware.ActorFrom(r.Context())
	if !ok {
		httpx.Error(w, apperrors.Unauthorized("não autenticado"))
		return
	}

	var input models.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.Error(w, apperrors.Validation("corpo da requisição inválido", nil))
		return
	}

	user, err := h.admin.CreateUser(r.Context(), actor, input)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, user)
}
