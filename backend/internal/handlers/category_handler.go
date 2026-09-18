package handlers

import (
	"encoding/json"
	"net/http"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/httpx"
	"ekaizen-backend/internal/middleware"
	"ekaizen-backend/internal/models"
	"ekaizen-backend/internal/repository"
	"ekaizen-backend/internal/services"
)

type CategoryHandler struct {
	categories repository.CategoryRepository
	admin      *services.AdminService
}

func NewCategoryHandler(categories repository.CategoryRepository, admin *services.AdminService) *CategoryHandler {
	return &CategoryHandler{categories: categories, admin: admin}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categories.List(r.Context())
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, categories)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	actor, ok := middleware.ActorFrom(r.Context())
	if !ok {
		httpx.Error(w, apperrors.Unauthorized("não autenticado"))
		return
	}

	var input models.CreateCategoryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.Error(w, apperrors.Validation("corpo da requisição inválido", nil))
		return
	}

	category, err := h.admin.CreateCategory(r.Context(), actor, input)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var input models.UpdateCategoryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.Error(w, apperrors.Validation("corpo da requisição inválido", nil))
		return
	}

	category, err := h.admin.UpdateCategory(r.Context(), actor, id, input)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, category)
}

func (h *CategoryHandler) SetApprovers(w http.ResponseWriter, r *http.Request) {
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

	var input models.SetApproversInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.Error(w, apperrors.Validation("corpo da requisição inválido", nil))
		return
	}

	category, err := h.admin.SetCategoryApprovers(r.Context(), actor, id, input)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, category)
}
