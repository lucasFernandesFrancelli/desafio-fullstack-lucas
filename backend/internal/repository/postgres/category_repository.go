package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/models"
)

type CategoryRepository struct {
	store *Store
}

func NewCategoryRepository(store *Store) *CategoryRepository {
	return &CategoryRepository{store: store}
}

func (r *CategoryRepository) List(ctx context.Context) ([]models.Category, error) {
	rows, err := r.store.db(ctx).Query(ctx, `
		SELECT c.id, c.name, COALESCE(c.description, ''),
		       ca.user_id, u.name, ca.approval_order
		FROM categories c
		LEFT JOIN category_approvers ca ON ca.category_id = c.id
		LEFT JOIN users u ON u.id = ca.user_id
		ORDER BY c.name, ca.approval_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byID := map[uuid.UUID]*models.Category{}
	var order []uuid.UUID

	for rows.Next() {
		var c models.Category
		var approverID *uuid.UUID
		var approverName *string
		var approvalOrder *int16

		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &approverID, &approverName, &approvalOrder); err != nil {
			return nil, err
		}

		cat, ok := byID[c.ID]
		if !ok {
			cat = &models.Category{ID: c.ID, Name: c.Name, Description: c.Description, Approvers: []models.CategoryApprover{}}
			byID[c.ID] = cat
			order = append(order, c.ID)
		}
		if approverID != nil {
			cat.Approvers = append(cat.Approvers, models.CategoryApprover{
				UserID: *approverID, UserName: *approverName, Order: *approvalOrder,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]models.Category, 0, len(order))
	for _, id := range order {
		result = append(result, *byID[id])
	}
	return result, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	categories, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, c := range categories {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, apperrors.NotFound("categoria não encontrada")
}

func (r *CategoryRepository) GetApproverByOrder(ctx context.Context, categoryID uuid.UUID, order int16) (*models.User, error) {
	row := r.store.db(ctx).QueryRow(ctx, `
		SELECT u.id, u.name, u.email, u.role, u.created_at
		FROM category_approvers ca
		JOIN users u ON u.id = ca.user_id
		WHERE ca.category_id = $1 AND ca.approval_order = $2`, categoryID, order)

	var u models.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound("aprovador não configurado para esta categoria/ordem")
		}
		return nil, err
	}
	return &u, nil
}

func (r *CategoryRepository) GetApproverOrder(ctx context.Context, categoryID, userID uuid.UUID) (*int16, error) {
	row := r.store.db(ctx).QueryRow(ctx, `
		SELECT approval_order FROM category_approvers
		WHERE category_id = $1 AND user_id = $2`, categoryID, userID)

	var order int16
	if err := row.Scan(&order); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}
