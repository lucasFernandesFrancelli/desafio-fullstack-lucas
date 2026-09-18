package handlers

import (
	"net/http"

	"ekaizen-backend/internal/httpx"
)

func Health(w http.ResponseWriter, r *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
