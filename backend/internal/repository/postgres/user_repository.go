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

type UserRepository struct {
	store *Store
}

func NewUserRepository(store *Store) *UserRepository {
	return &UserRepository{store: store}
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	row := r.store.db(ctx).QueryRow(ctx,
		`SELECT id, name, email, role, created_at FROM users WHERE id = $1`, id)

	var u models.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound("usuário não encontrado")
		}
		return nil, err
	}
	return &u, nil
}

// ListProfiles retorna todos os usuários com as categorias em que cada um é
// aprovador, para alimentar a tela de seletor de perfil.
func (r *UserRepository) ListProfiles(ctx context.Context) ([]models.UserProfile, error) {
	rows, err := r.store.db(ctx).Query(ctx, `
		SELECT u.id, u.name, u.email, u.role, u.created_at,
		       ca.category_id, c.name, ca.approval_order
		FROM users u
		LEFT JOIN category_approvers ca ON ca.user_id = u.id
		LEFT JOIN categories c ON c.id = ca.category_id
		ORDER BY u.name, ca.approval_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byID := map[uuid.UUID]*models.UserProfile{}
	var order []uuid.UUID

	for rows.Next() {
		var u models.User
		var catID *uuid.UUID
		var catName *string
		var approvalOrder *int16

		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt,
			&catID, &catName, &approvalOrder); err != nil {
			return nil, err
		}

		profile, ok := byID[u.ID]
		if !ok {
			profile = &models.UserProfile{User: u, ApproverFor: []models.ApproverFor{}}
			byID[u.ID] = profile
			order = append(order, u.ID)
		}
		if catID != nil {
			profile.ApproverFor = append(profile.ApproverFor, models.ApproverFor{
				CategoryID:   *catID,
				CategoryName: *catName,
				Order:        *approvalOrder,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]models.UserProfile, 0, len(order))
	for _, id := range order {
		result = append(result, *byID[id])
	}
	return result, nil
}

// List retorna todos os usuários, para alimentar a área de gestão de
// cadastros (selects de aprovador, listagem de pessoas).
func (r *UserRepository) List(ctx context.Context) ([]models.User, error) {
	rows, err := r.store.db(ctx).Query(ctx,
		`SELECT id, name, email, role, created_at FROM users ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, rows.Err()
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	row := r.store.db(ctx).QueryRow(ctx, `
		INSERT INTO users (id, name, email, role) VALUES ($1,$2,$3,$4)
		RETURNING created_at`, u.ID, u.Name, u.Email, u.Role)

	if err := row.Scan(&u.CreatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperrors.Validation("e-mail já cadastrado", map[string]string{"email": "já está em uso"})
		}
		return err
	}
	return nil
}
