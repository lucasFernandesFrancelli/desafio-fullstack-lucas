package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/models"
	"ekaizen-backend/internal/services"
	"ekaizen-backend/internal/testutil"
)

type adminHarness struct {
	users *testutil.UserRepo
	cats  *testutil.CategoryRepo
	svc   *services.AdminService
}

func newAdminHarness() *adminHarness {
	h := &adminHarness{users: &testutil.UserRepo{}, cats: &testutil.CategoryRepo{}}
	h.svc = services.NewAdminService(&testutil.Transactor{}, h.users, h.cats)
	return h
}

var gestor = models.User{ID: uuid.New(), Role: models.RoleGestor}
var colaborador = models.User{ID: uuid.New(), Role: models.RoleColaborador}

func TestAdminListUsers_ForbiddenForNonGestor(t *testing.T) {
	h := newAdminHarness()
	_, err := h.svc.ListUsers(context.Background(), colaborador)
	requireAppError(t, err, apperrors.CodeForbidden)
}

func TestAdminListUsers_Success(t *testing.T) {
	h := newAdminHarness()
	h.users.On("List", mock.Anything).Return([]models.User{{ID: uuid.New(), Name: "Fulano"}}, nil)

	users, err := h.svc.ListUsers(context.Background(), gestor)
	require.NoError(t, err)
	require.Len(t, users, 1)
}

func TestAdminCreateUser_ValidationMissingFields(t *testing.T) {
	h := newAdminHarness()
	_, err := h.svc.CreateUser(context.Background(), gestor, models.CreateUserInput{})
	requireAppError(t, err, apperrors.CodeValidation)
	appErr, _ := apperrors.As(err)
	require.Contains(t, appErr.Fields, "name")
	require.Contains(t, appErr.Fields, "email")
	require.Contains(t, appErr.Fields, "role")
}

func TestAdminCreateUser_ForbiddenForNonGestor(t *testing.T) {
	h := newAdminHarness()
	_, err := h.svc.CreateUser(context.Background(), colaborador, models.CreateUserInput{
		Name: "Fulano", Email: "fulano@ekaizen.example", Role: models.RoleColaborador,
	})
	requireAppError(t, err, apperrors.CodeForbidden)
}

func TestAdminCreateUser_Success(t *testing.T) {
	h := newAdminHarness()
	h.users.On("Create", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
		return u.Name == "Nova Pessoa" && u.Email == "nova@ekaizen.example" && u.Role == models.RoleAnalista
	})).Return(nil)

	user, err := h.svc.CreateUser(context.Background(), gestor, models.CreateUserInput{
		Name: "Nova Pessoa", Email: "nova@ekaizen.example", Role: models.RoleAnalista,
	})
	require.NoError(t, err)
	require.Equal(t, "Nova Pessoa", user.Name)
}

func TestAdminCreateCategory_ValidationMissingFields(t *testing.T) {
	h := newAdminHarness()
	_, err := h.svc.CreateCategory(context.Background(), gestor, models.CreateCategoryInput{})
	requireAppError(t, err, apperrors.CodeValidation)
	appErr, _ := apperrors.As(err)
	require.Contains(t, appErr.Fields, "name")
	require.Contains(t, appErr.Fields, "firstApproverId")
	require.Contains(t, appErr.Fields, "secondApproverId")
}

func TestAdminCreateCategory_SameApproverTwice(t *testing.T) {
	h := newAdminHarness()
	sameID := uuid.New()
	_, err := h.svc.CreateCategory(context.Background(), gestor, models.CreateCategoryInput{
		Name: "Ergonomia", FirstApproverID: sameID, SecondApproverID: sameID,
	})
	requireAppError(t, err, apperrors.CodeValidation)
	appErr, _ := apperrors.As(err)
	require.Contains(t, appErr.Fields, "secondApproverId")
}

func TestAdminCreateCategory_ApproverNotFound(t *testing.T) {
	h := newAdminHarness()
	firstID, secondID := uuid.New(), uuid.New()
	h.users.On("GetByID", mock.Anything, firstID).Return(nil, apperrors.NotFound("não encontrado"))

	_, err := h.svc.CreateCategory(context.Background(), gestor, models.CreateCategoryInput{
		Name: "Ergonomia", FirstApproverID: firstID, SecondApproverID: secondID,
	})
	requireAppError(t, err, apperrors.CodeValidation)
}

func TestAdminCreateCategory_Success(t *testing.T) {
	h := newAdminHarness()
	firstID, secondID := uuid.New(), uuid.New()
	created := &models.Category{ID: uuid.New(), Name: "Ergonomia"}

	h.users.On("GetByID", mock.Anything, firstID).Return(&models.User{ID: firstID}, nil)
	h.users.On("GetByID", mock.Anything, secondID).Return(&models.User{ID: secondID}, nil)
	h.cats.On("Create", mock.Anything, mock.MatchedBy(func(c *models.Category) bool { return c.Name == "Ergonomia" })).
		Run(func(args mock.Arguments) {
			c := args.Get(1).(*models.Category)
			created.ID = c.ID
		}).Return(nil)
	h.cats.On("SetApprovers", mock.Anything, mock.Anything, firstID, secondID).Return(nil)
	h.cats.On("GetByID", mock.Anything, mock.Anything).Return(created, nil)

	category, err := h.svc.CreateCategory(context.Background(), gestor, models.CreateCategoryInput{
		Name: "Ergonomia", FirstApproverID: firstID, SecondApproverID: secondID,
	})
	require.NoError(t, err)
	require.Equal(t, "Ergonomia", category.Name)
}

func TestAdminUpdateCategory_Success(t *testing.T) {
	h := newAdminHarness()
	id := uuid.New()
	existing := &models.Category{ID: id, Name: "Antigo", Description: "desc antiga"}
	updated := &models.Category{ID: id, Name: "Novo Nome", Description: "desc antiga"}

	h.cats.On("GetByID", mock.Anything, id).Return(existing, nil).Once()
	newName := "Novo Nome"
	h.cats.On("UpdateInfo", mock.Anything, id, "Novo Nome", "desc antiga").Return(nil)
	h.cats.On("GetByID", mock.Anything, id).Return(updated, nil).Once()

	category, err := h.svc.UpdateCategory(context.Background(), gestor, id, models.UpdateCategoryInput{Name: &newName})
	require.NoError(t, err)
	require.Equal(t, "Novo Nome", category.Name)
}

func TestAdminSetCategoryApprovers_SameApproverTwice(t *testing.T) {
	h := newAdminHarness()
	sameID := uuid.New()
	_, err := h.svc.SetCategoryApprovers(context.Background(), gestor, uuid.New(), models.SetApproversInput{
		FirstApproverID: sameID, SecondApproverID: sameID,
	})
	requireAppError(t, err, apperrors.CodeValidation)
}

func TestAdminSetCategoryApprovers_Success(t *testing.T) {
	h := newAdminHarness()
	categoryID := uuid.New()
	firstID, secondID := uuid.New(), uuid.New()
	category := &models.Category{ID: categoryID, Name: "Qualidade"}

	h.users.On("GetByID", mock.Anything, firstID).Return(&models.User{ID: firstID}, nil)
	h.users.On("GetByID", mock.Anything, secondID).Return(&models.User{ID: secondID}, nil)
	h.cats.On("GetByID", mock.Anything, categoryID).Return(category, nil)
	h.cats.On("SetApprovers", mock.Anything, categoryID, firstID, secondID).Return(nil)

	result, err := h.svc.SetCategoryApprovers(context.Background(), gestor, categoryID, models.SetApproversInput{
		FirstApproverID: firstID, SecondApproverID: secondID,
	})
	require.NoError(t, err)
	require.Equal(t, categoryID, result.ID)
}

func TestAdminSetCategoryApprovers_ForbiddenForNonGestor(t *testing.T) {
	h := newAdminHarness()
	_, err := h.svc.SetCategoryApprovers(context.Background(), colaborador, uuid.New(), models.SetApproversInput{
		FirstApproverID: uuid.New(), SecondApproverID: uuid.New(),
	})
	requireAppError(t, err, apperrors.CodeForbidden)
}
