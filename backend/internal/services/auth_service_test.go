package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/auth"
	"ekaizen-backend/internal/models"
	"ekaizen-backend/internal/services"
	"ekaizen-backend/internal/testutil"
)

func TestAuthService_Login_Success(t *testing.T) {
	users := &testutil.UserRepo{}
	issuer := auth.NewJWTIssuer("test-secret")
	svc := services.NewAuthService(users, issuer)

	user := &models.User{ID: uuid.New(), Name: "Ana", Role: models.RoleColaborador}
	users.On("GetByID", mock.Anything, user.ID).Return(user, nil)

	token, got, err := svc.Login(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, user.ID, got.ID)

	claims, err := issuer.Parse(token)
	require.NoError(t, err)
	require.Equal(t, user.ID, claims.UserID)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	users := &testutil.UserRepo{}
	issuer := auth.NewJWTIssuer("test-secret")
	svc := services.NewAuthService(users, issuer)

	missingID := uuid.New()
	users.On("GetByID", mock.Anything, missingID).Return(nil, apperrors.NotFound("usuário não encontrado"))

	_, _, err := svc.Login(context.Background(), missingID)
	appErr, ok := apperrors.As(err)
	require.True(t, ok)
	require.Equal(t, apperrors.CodeNotFound, appErr.Code)
}

func TestAuthService_Me(t *testing.T) {
	users := &testutil.UserRepo{}
	issuer := auth.NewJWTIssuer("test-secret")
	svc := services.NewAuthService(users, issuer)

	user := &models.User{ID: uuid.New(), Name: "Ana"}
	users.On("GetByID", mock.Anything, user.ID).Return(user, nil)

	got, err := svc.Me(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, user.Name, got.Name)
}

func TestAuthService_ListProfiles(t *testing.T) {
	users := &testutil.UserRepo{}
	issuer := auth.NewJWTIssuer("test-secret")
	svc := services.NewAuthService(users, issuer)

	profiles := []models.UserProfile{{User: models.User{Name: "Ana"}}}
	users.On("ListProfiles", mock.Anything).Return(profiles, nil)

	got, err := svc.ListProfiles(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
}
