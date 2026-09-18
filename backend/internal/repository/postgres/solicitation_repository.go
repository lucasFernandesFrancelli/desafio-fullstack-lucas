package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"ekaizen-backend/internal/apperrors"
	"ekaizen-backend/internal/models"
)

type SolicitationRepository struct {
	store *Store
}

func NewSolicitationRepository(store *Store) *SolicitationRepository {
	return &SolicitationRepository{store: store}
}

const solicitationCoreColumns = `
	id, title, problem_description, proposed_improvement, category_id, location,
	requester_id, status, current_approval_step, severity, urgency, trend, priority,
	analysis_notes, rejection_reason, last_transition_at, created_at, updated_at`

func scanSolicitationCore(row pgx.Row, s *models.Solicitation) error {
	return row.Scan(
		&s.ID, &s.Title, &s.ProblemDescription, &s.ProposedImprovement, &s.CategoryID, &s.Location,
		&s.RequesterID, &s.Status, &s.CurrentApprovalStep, &s.Severity, &s.Urgency, &s.Trend, &s.Priority,
		&s.AnalysisNotes, &s.RejectionReason, &s.LastTransitionAt, &s.CreatedAt, &s.UpdatedAt,
	)
}

func (r *SolicitationRepository) Create(ctx context.Context, s *models.Solicitation) error {
	row := r.store.db(ctx).QueryRow(ctx, `
		INSERT INTO solicitations (id, title, problem_description, proposed_improvement, category_id, location, requester_id, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING created_at, updated_at, last_transition_at`,
		s.ID, s.Title, s.ProblemDescription, s.ProposedImprovement, s.CategoryID, s.Location, s.RequesterID, s.Status)
	return row.Scan(&s.CreatedAt, &s.UpdatedAt, &s.LastTransitionAt)
}

func (r *SolicitationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Solicitation, error) {
	row := r.store.db(ctx).QueryRow(ctx,
		fmt.Sprintf(`SELECT %s FROM solicitations WHERE id = $1`, solicitationCoreColumns), id)
	var s models.Solicitation
	if err := scanSolicitationCore(row, &s); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound("solicitação não encontrada")
		}
		return nil, err
	}
	if err := r.attachNames(ctx, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// GetForUpdate trava a linha (FOR UPDATE) para o restante da transação —
// deve ser chamado apenas dentro de Store.WithinTx.
func (r *SolicitationRepository) GetForUpdate(ctx context.Context, id uuid.UUID) (*models.Solicitation, error) {
	row := r.store.db(ctx).QueryRow(ctx,
		fmt.Sprintf(`SELECT %s FROM solicitations WHERE id = $1 FOR UPDATE`, solicitationCoreColumns), id)
	var s models.Solicitation
	if err := scanSolicitationCore(row, &s); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound("solicitação não encontrada")
		}
		return nil, err
	}
	return &s, nil
}

func (r *SolicitationRepository) attachNames(ctx context.Context, s *models.Solicitation) error {
	row := r.store.db(ctx).QueryRow(ctx, `SELECT name FROM users WHERE id = $1`, s.RequesterID)
	if err := row.Scan(&s.RequesterName); err != nil {
		return err
	}
	if s.CategoryID != nil {
		row := r.store.db(ctx).QueryRow(ctx, `SELECT name FROM categories WHERE id = $1`, *s.CategoryID)
		if err := row.Scan(&s.CategoryName); err != nil {
			return err
		}
	}
	return nil
}

func (r *SolicitationRepository) Update(ctx context.Context, s *models.Solicitation) error {
	row := r.store.db(ctx).QueryRow(ctx, `
		UPDATE solicitations SET
			title = $2, problem_description = $3, proposed_improvement = $4, category_id = $5,
			location = $6, status = $7, current_approval_step = $8, severity = $9, urgency = $10,
			trend = $11, analysis_notes = $12, rejection_reason = $13, last_transition_at = $14,
			updated_at = now()
		WHERE id = $1
		RETURNING updated_at, priority`,
		s.ID, s.Title, s.ProblemDescription, s.ProposedImprovement, s.CategoryID,
		s.Location, s.Status, s.CurrentApprovalStep, s.Severity, s.Urgency,
		s.Trend, s.AnalysisNotes, s.RejectionReason, s.LastTransitionAt)
	return row.Scan(&s.UpdatedAt, &s.Priority)
}

func (r *SolicitationRepository) List(ctx context.Context, filter models.SolicitationFilter) ([]models.Solicitation, error) {
	query := fmt.Sprintf(`SELECT %s FROM solicitations WHERE 1=1`, solicitationCoreColumns)
	args := []any{}
	argN := 1

	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argN)
		args = append(args, *filter.Status)
		argN++
	}
	if filter.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", argN)
		args = append(args, *filter.CategoryID)
		argN++
	}
	if strings.TrimSpace(filter.Query) != "" {
		query += fmt.Sprintf(` AND (title ILIKE $%d OR problem_description ILIKE $%d OR proposed_improvement ILIKE $%d)`, argN, argN, argN)
		args = append(args, "%"+strings.TrimSpace(filter.Query)+"%")
		argN++
	}
	query += " ORDER BY updated_at DESC LIMIT 200"

	rows, err := r.store.db(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Solicitation
	for rows.Next() {
		var s models.Solicitation
		if err := scanSolicitationCore(rows, &s); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range result {
		if err := r.attachNames(ctx, &result[i]); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (r *SolicitationRepository) CountsByStatus(ctx context.Context) (map[models.Status]int, error) {
	rows, err := r.store.db(ctx).Query(ctx, `SELECT status, COUNT(*) FROM solicitations GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[models.Status]int{
		models.StatusRascunho:    0,
		models.StatusEmAprovacao: 0,
		models.StatusEmAnalise:   0,
		models.StatusFinalizada:  0,
		models.StatusRecusada:    0,
	}
	for rows.Next() {
		var status models.Status
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		result[status] = count
	}
	return result, rows.Err()
}

func (r *SolicitationRepository) OldestPending(ctx context.Context, limit int) ([]models.OldestPendingItem, error) {
	rows, err := r.store.db(ctx).Query(ctx, `
		SELECT s.id, s.title, s.status, COALESCE(c.name, ''),
		       COALESCE(
		         CASE
		           WHEN s.status = 'em_aprovacao' THEN (
		             SELECT u.name FROM category_approvers ca
		             JOIN users u ON u.id = ca.user_id
		             WHERE ca.category_id = s.category_id AND ca.approval_order = s.current_approval_step
		           )
		           WHEN s.status = 'em_analise' THEN 'Qualquer analista'
		           WHEN s.status = 'rascunho' THEN (SELECT name FROM users WHERE id = s.requester_id)
		           ELSE ''
		         END, ''),
		       EXTRACT(DAY FROM now() - s.last_transition_at)::int
		FROM solicitations s
		LEFT JOIN categories c ON c.id = s.category_id
		WHERE s.status NOT IN ('finalizada', 'recusada')
		ORDER BY s.last_transition_at ASC
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.OldestPendingItem
	for rows.Next() {
		var item models.OldestPendingItem
		if err := rows.Scan(&item.ID, &item.Title, &item.Status, &item.CategoryName,
			&item.PendingActorName, &item.DaysSinceLastMove); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *SolicitationRepository) AwaitingApprovalByUser(ctx context.Context) ([]models.UserPendingCount, error) {
	rows, err := r.store.db(ctx).Query(ctx, `
		SELECT u.id, u.name, COUNT(*)
		FROM solicitations s
		JOIN category_approvers ca ON ca.category_id = s.category_id AND ca.approval_order = s.current_approval_step
		JOIN users u ON u.id = ca.user_id
		WHERE s.status = 'em_aprovacao'
		GROUP BY u.id, u.name
		ORDER BY COUNT(*) DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.UserPendingCount
	for rows.Next() {
		var item models.UserPendingCount
		if err := rows.Scan(&item.UserID, &item.UserName, &item.Count); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
