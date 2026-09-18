package postgres

import (
	"context"

	"github.com/google/uuid"

	"ekaizen-backend/internal/models"
)

type HistoryRepository struct {
	store *Store
}

func NewHistoryRepository(store *Store) *HistoryRepository {
	return &HistoryRepository{store: store}
}

func (r *HistoryRepository) Create(ctx context.Context, h *models.HistoryEntry) error {
	row := r.store.db(ctx).QueryRow(ctx, `
		INSERT INTO solicitation_history
			(id, solicitation_id, from_status, to_status, action, actor_id, approval_step, comment, severity, urgency, trend)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING created_at`,
		h.ID, h.SolicitationID, h.FromStatus, h.ToStatus, h.Action, h.ActorID,
		h.ApprovalStep, h.Comment, h.Severity, h.Urgency, h.Trend)
	return row.Scan(&h.CreatedAt)
}

func (r *HistoryRepository) ListBySolicitation(ctx context.Context, solicitationID uuid.UUID) ([]models.HistoryEntry, error) {
	rows, err := r.store.db(ctx).Query(ctx, `
		SELECT h.id, h.from_status, h.to_status, h.action, h.actor_id, u.name,
		       h.approval_step, h.comment, h.severity, h.urgency, h.trend, h.created_at
		FROM solicitation_history h
		JOIN users u ON u.id = h.actor_id
		WHERE h.solicitation_id = $1
		ORDER BY h.created_at ASC`, solicitationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.HistoryEntry
	for rows.Next() {
		var h models.HistoryEntry
		if err := rows.Scan(&h.ID, &h.FromStatus, &h.ToStatus, &h.Action, &h.ActorID, &h.ActorName,
			&h.ApprovalStep, &h.Comment, &h.Severity, &h.Urgency, &h.Trend, &h.CreatedAt); err != nil {
			return nil, err
		}
		h.SolicitationID = solicitationID
		result = append(result, h)
	}
	return result, rows.Err()
}
