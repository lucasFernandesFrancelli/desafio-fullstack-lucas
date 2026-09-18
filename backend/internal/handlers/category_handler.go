package handlers

import (
	"net/http"

	"ekaizen-backend/internal/httpx"
	"ekaizen-backend/internal/repository"
)

type CategoryHandler struct {
	categories repository.CategoryRepository
}

func NewCategoryHandler(categories repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{categories: categories}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categories.List(r.Context())
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, categories)
}
