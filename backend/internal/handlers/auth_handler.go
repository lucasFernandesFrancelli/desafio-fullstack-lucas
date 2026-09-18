package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/httpx"
	"ekaizen-backend/internal/middleware"
	"ekaizen-backend/internal/services"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) ListProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.service.ListProfiles(r.Context())
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, profiles)
}

type loginRequest struct {
	UserID string `json:"userId"`
}

type loginResponse struct {
	Token string `json:"token"`
	User  any    `json:"user"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, apperrors.Validation("corpo da requisição inválido", nil))
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		httpx.Error(w, apperrors.Validation("userId inválido", map[string]string{"userId": "deve ser um uuid válido"}))
		return
	}

	token, user, err := h.service.Login(r.Context(), userID)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, loginResponse{Token: token, User: user})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	actor, ok := middleware.ActorFrom(r.Context())
	if !ok {
		httpx.Error(w, apperrors.Unauthorized("não autenticado"))
		return
	}
	user, err := h.service.Me(r.Context(), actor.ID)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, user)
}
