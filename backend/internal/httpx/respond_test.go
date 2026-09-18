package httpx_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/httpx"
)

func TestJSON_WritesStatusAndBody(t *testing.T) {
	rec := httptest.NewRecorder()
	httpx.JSON(rec, http.StatusCreated, map[string]string{"ok": "sim"})

	assert.Equal(t, http.StatusCreated, rec.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "sim", body["ok"])
}

func TestError_MapsAppErrorCodes(t *testing.T) {
	cases := []struct {
		err            error
		expectedStatus int
	}{
		{apperrors.NotFound("x"), http.StatusNotFound},
		{apperrors.Forbidden("x"), http.StatusForbidden},
		{apperrors.Conflict("x"), http.StatusConflict},
		{apperrors.Validation("x", nil), http.StatusUnprocessableEntity},
		{apperrors.Unauthorized("x"), http.StatusUnauthorized},
	}
	for _, c := range cases {
		rec := httptest.NewRecorder()
		httpx.Error(rec, c.err)
		assert.Equal(t, c.expectedStatus, rec.Code)
	}
}

func TestError_UnknownErrorBecomes500(t *testing.T) {
	rec := httptest.NewRecorder()
	httpx.Error(rec, errors.New("algo inesperado"))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var body map[string]map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "internal_error", body["error"]["code"])
}
