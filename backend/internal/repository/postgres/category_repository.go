package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

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

func (r *CategoryRepository) Create(ctx context.Context, c *models.Category) error {
	_, err := r.store.db(ctx).Exec(ctx,
		`INSERT INTO categories (id, name, description) VALUES ($1,$2,$3)`, c.ID, c.Name, c.Description)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperrors.Validation("categoria já existe", map[string]string{"name": "já está em uso"})
		}
		return err
	}
	return nil
}

func (r *CategoryRepository) UpdateInfo(ctx context.Context, id uuid.UUID, name, description string) error {
	_, err := r.store.db(ctx).Exec(ctx,
		`UPDATE categories SET name = $2, description = $3 WHERE id = $1`, id, name, description)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperrors.Validation("categoria já existe", map[string]string{"name": "já está em uso"})
		}
		return err
	}
	return nil
}

// SetApprovers substitui os dois aprovadores de uma categoria de uma vez só
// (delete + insert). Fazer isso em dois passos evita violar a constraint
// UNIQUE(category_id, user_id) em trocas cruzadas (ex.: A vira 2º aprovador
// e B vira 1º, quando antes era o contrário).
func (r *CategoryRepository) SetApprovers(ctx context.Context, categoryID, firstApproverID, secondApproverID uuid.UUID) error {
	if _, err := r.store.db(ctx).Exec(ctx,
		`DELETE FROM category_approvers WHERE category_id = $1`, categoryID); err != nil {
		return err
	}
	_, err := r.store.db(ctx).Exec(ctx, `
		INSERT INTO category_approvers (id, category_id, user_id, approval_order)
		VALUES ($1,$2,$3,1), ($4,$2,$5,2)`,
		uuid.New(), categoryID, firstApproverID, uuid.New(), secondApproverID)
	return err
}
