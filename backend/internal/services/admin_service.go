package services

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/models"
	"ekaizen-backend/internal/repository"
)

// AdminService cobre a área de gestão de cadastros: quem são os usuários e
// quem aprova cada categoria. Existe para que o avaliador consiga conferir
// (e alterar) essa configuração pela própria aplicação, em vez de enxergar
// só dados fixos no seed. Toda operação aqui é restrita ao papel "gestor".
type AdminService struct {
	tx         repository.Transactor
	users      repository.UserRepository
	categories repository.CategoryRepository
}

func NewAdminService(tx repository.Transactor, users repository.UserRepository, categories repository.CategoryRepository) *AdminService {
	return &AdminService{tx: tx, users: users, categories: categories}
}

func requireGestor(actor models.User) error {
	if actor.Role != models.RoleGestor {
		return apperrors.Forbidden("apenas o gestor pode gerenciar cadastros")
	}
	return nil
}

func (s *AdminService) ListUsers(ctx context.Context, actor models.User) ([]models.User, error) {
	if err := requireGestor(actor); err != nil {
		return nil, err
	}
	return s.users.List(ctx)
}

func (s *AdminService) CreateUser(ctx context.Context, actor models.User, input models.CreateUserInput) (*models.User, error) {
	if err := requireGestor(actor); err != nil {
		return nil, err
	}

	name := strings.TrimSpace(input.Name)
	email := strings.TrimSpace(input.Email)

	fields := map[string]string{}
	if name == "" {
		fields["name"] = "obrigatório"
	}
	if email == "" {
		fields["email"] = "obrigatório"
	}
	if !input.Role.Valid() {
		fields["role"] = "deve ser colaborador, analista ou gestor"
	}
	if len(fields) > 0 {
		return nil, apperrors.Validation("preencha os campos obrigatórios", fields)
	}

	user := &models.User{ID: uuid.New(), Name: name, Email: email, Role: input.Role}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AdminService) CreateCategory(ctx context.Context, actor models.User, input models.CreateCategoryInput) (*models.Category, error) {
	if err := requireGestor(actor); err != nil {
		return nil, err
	}

	name := strings.TrimSpace(input.Name)
	fields := map[string]string{}
	if name == "" {
		fields["name"] = "obrigatório"
	}
	if input.FirstApproverID == uuid.Nil {
		fields["firstApproverId"] = "obrigatório"
	}
	if input.SecondApproverID == uuid.Nil {
		fields["secondApproverId"] = "obrigatório"
	}
	if input.FirstApproverID != uuid.Nil && input.FirstApproverID == input.SecondApproverID {
		fields["secondApproverId"] = "deve ser uma pessoa diferente do 1º aprovador"
	}
	if len(fields) > 0 {
		return nil, apperrors.Validation("preencha os campos obrigatórios", fields)
	}

	if err := s.validateApprovers(ctx, input.FirstApproverID, input.SecondApproverID); err != nil {
		return nil, err
	}

	var categoryID uuid.UUID
	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		category := &models.Category{ID: uuid.New(), Name: name, Description: strings.TrimSpace(input.Description)}
		if err := s.categories.Create(ctx, category); err != nil {
			return err
		}
		categoryID = category.ID
		return s.categories.SetApprovers(ctx, categoryID, input.FirstApproverID, input.SecondApproverID)
	})
	if err != nil {
		return nil, err
	}

	return s.categories.GetByID(ctx, categoryID)
}

func (s *AdminService) UpdateCategory(ctx context.Context, actor models.User, id uuid.UUID, input models.UpdateCategoryInput) (*models.Category, error) {
	if err := requireGestor(actor); err != nil {
		return nil, err
	}

	existing, err := s.categories.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	name := existing.Name
	if input.Name != nil {
		name = strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, apperrors.Validation("nome é obrigatório", map[string]string{"name": "obrigatório"})
		}
	}
	description := existing.Description
	if input.Description != nil {
		description = strings.TrimSpace(*input.Description)
	}

	if err := s.categories.UpdateInfo(ctx, id, name, description); err != nil {
		return nil, err
	}
	return s.categories.GetByID(ctx, id)
}

func (s *AdminService) SetCategoryApprovers(ctx context.Context, actor models.User, id uuid.UUID, input models.SetApproversInput) (*models.Category, error) {
	if err := requireGestor(actor); err != nil {
		return nil, err
	}

	if input.FirstApproverID == input.SecondApproverID {
		return nil, apperrors.Validation("os dois aprovadores devem ser pessoas distintas",
			map[string]string{"secondApproverId": "deve ser diferente do 1º aprovador"})
	}
	if err := s.validateApprovers(ctx, input.FirstApproverID, input.SecondApproverID); err != nil {
		return nil, err
	}
	if _, err := s.categories.GetByID(ctx, id); err != nil {
		return nil, err
	}

	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		return s.categories.SetApprovers(ctx, id, input.FirstApproverID, input.SecondApproverID)
	})
	if err != nil {
		return nil, err
	}
	return s.categories.GetByID(ctx, id)
}

func (s *AdminService) validateApprovers(ctx context.Context, firstApproverID, secondApproverID uuid.UUID) error {
	if _, err := s.users.GetByID(ctx, firstApproverID); err != nil {
		return apperrors.Validation("1º aprovador inválido", map[string]string{"firstApproverId": "usuário não encontrado"})
	}
	if _, err := s.users.GetByID(ctx, secondApproverID); err != nil {
		return apperrors.Validation("2º aprovador inválido", map[string]string{"secondApproverId": "usuário não encontrado"})
	}
	return nil
}
