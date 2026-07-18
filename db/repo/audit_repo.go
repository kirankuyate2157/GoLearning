package repo

import (
	"context"

	"app/models"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditRepo struct {
	db *pgxpool.Pool
}

func NewAuditRepo(db *pgxpool.Pool) *AuditRepo {
	return &AuditRepo{db: db}
}

// LogAuditEvent records a critical action in the audit_events table
func (r *AuditRepo) LogAuditEvent(ctx context.Context, userUID, action, resourceType, resourceUID, oldValue, newValue, notes string) (*models.AuditEvent, error) {
	row := r.db.QueryRow(ctx, `INSERT INTO audit_events (user_uid, action, resource_type, resource_uid, old_value, new_value, notes) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING uid, user_uid, action, resource_type, resource_uid, old_value, new_value, notes, created_at`, userUID, action, resourceType, resourceUID, oldValue, newValue, notes)
	ae := &models.AuditEvent{}
	var uid pgtype.UUID
	var createdAt pgtype.Timestamptz
	if err := row.Scan(&uid, &ae.UserUID, &ae.Action, &ae.ResourceType, &ae.ResourceUID, &ae.OldValue, &ae.NewValue, &ae.Notes, &createdAt); err != nil {
		return nil, err
	}
	ae.UID = uid.String()
	if createdAt.Valid {
		t := createdAt.Time
		ae.CreatedAt = &t
	}
	return ae, nil
}

// ListAuditEventsByUser returns all audit events for a given user
func (r *AuditRepo) ListAuditEventsByUser(ctx context.Context, userUID string, limit int) ([]models.AuditEvent, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := r.db.Query(ctx, `SELECT uid, user_uid, action, resource_type, resource_uid, old_value, new_value, notes, created_at FROM audit_events WHERE user_uid = $1 ORDER BY created_at DESC LIMIT $2`, userUID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.AuditEvent
	for rows.Next() {
		var ae models.AuditEvent
		var uid pgtype.UUID
		var createdAt pgtype.Timestamptz
		if err := rows.Scan(&uid, &ae.UserUID, &ae.Action, &ae.ResourceType, &ae.ResourceUID, &ae.OldValue, &ae.NewValue, &ae.Notes, &createdAt); err != nil {
			return nil, err
		}
		ae.UID = uid.String()
		if createdAt.Valid {
			t := createdAt.Time
			ae.CreatedAt = &t
		}
		out = append(out, ae)
	}
	return out, nil
}

// ListAuditEventsByResource returns all audit events for a specific resource
func (r *AuditRepo) ListAuditEventsByResource(ctx context.Context, resourceType, resourceUID string, limit int) ([]models.AuditEvent, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := r.db.Query(ctx, `SELECT uid, user_uid, action, resource_type, resource_uid, old_value, new_value, notes, created_at FROM audit_events WHERE resource_type = $1 AND resource_uid = $2 ORDER BY created_at DESC LIMIT $3`, resourceType, resourceUID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.AuditEvent
	for rows.Next() {
		var ae models.AuditEvent
		var uid pgtype.UUID
		var createdAt pgtype.Timestamptz
		if err := rows.Scan(&uid, &ae.UserUID, &ae.Action, &ae.ResourceType, &ae.ResourceUID, &ae.OldValue, &ae.NewValue, &ae.Notes, &createdAt); err != nil {
			return nil, err
		}
		ae.UID = uid.String()
		if createdAt.Valid {
			t := createdAt.Time
			ae.CreatedAt = &t
		}
		out = append(out, ae)
	}
	return out, nil
}
