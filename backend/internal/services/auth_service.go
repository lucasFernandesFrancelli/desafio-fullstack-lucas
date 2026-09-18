package services

import (
	"context"

	"github.com/google/uuid"

	"ekaizen-backend/internal/auth"
	"ekaizen-backend/internal/models"
	"ekaizen-backend/internal/repository"
)

// AuthService implementa o "login" mock do desafio: não há senha, o
// frontend lista os perfis disponíveis e o usuário escolhe um. O backend
// ainda assim emite um JWT real para autenticar as próximas chamadas.
type AuthService struct {
	users  repository.UserRepository
	issuer *auth.JWTIssuer
}

func NewAuthService(users repository.UserRepository, issuer *auth.JWTIssuer) *AuthService {
	return &AuthService{users: users, issuer: issuer}
}

func (s *AuthService) ListProfiles(ctx context.Context) ([]models.UserProfile, error) {
	return s.users.ListProfiles(ctx)
}

func (s *AuthService) Login(ctx context.Context, userID uuid.UUID) (string, *models.User, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return "", nil, err
	}
	token, err := s.issuer.Issue(*user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *AuthService) Me(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.users.GetByID(ctx, userID)
}
