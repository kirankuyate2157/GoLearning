package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"app/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WarrantyRow represents a warranty row from the database
type WarrantyRow struct {
	ID          int
	UID         pgtype.Text
	Type        string
	RefID       int
	Provider    pgtype.Int8
	StartDate   pgtype.Timestamp
	EndDate     pgtype.Timestamp
	CreatedAt   pgtype.Timestamp
	CreatedBy   pgtype.Int8
	RefUID      pgtype.Text
	ProviderUID pgtype.Text
}

type WarrantyPatchReq struct {
	Provider  *int
	StartDate *time.Time
	EndDate   *time.Time
}

// GetWarranties fetches warranties for a specific type and ref_id
func GetWarranties(pool *pgxpool.Pool, qp models.Warranty) ([]WarrantyRow, error) {
	ctx := context.Background()
	var filters []string
	var filterValues []any
	idx := 1
	query := `
		SELECT 
			w.id,
			w.uid,
			w.type,
			w.ref_id,
			w.provider,
			w.start_date,
			w.end_date,
			w.created_at,
			w.created_by,
			m.uid ref_uid,
			v.uid vendor_uid
		FROM warranties w
		left join machines m on w.ref_id = m.id
		left join vendors v on w.provider = v.id
		
	`

	if len(qp.RefUID) != 0 {
		filters = append(filters, fmt.Sprintf(" m.uid = $%d ", idx))
		filterValues = append(filterValues, qp.RefUID)
		idx++
	}

	if len(qp.ProviderUID) != 0 {
		filters = append(filters, fmt.Sprintf(" v.uid = $%d ", idx))
		filterValues = append(filterValues, qp.ProviderUID)
		idx++
	}

	if len(qp.UID) != 0 {
		filters = append(filters, fmt.Sprintf(" w.uid = $%d ", idx))
		filterValues = append(filterValues, qp.UID)
		idx++
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += " ORDER BY w.id ASC"

	rows, err := pool.Query(ctx, query, filterValues...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var out []WarrantyRow
	for rows.Next() {
		var r WarrantyRow
		if err := rows.Scan(&r.ID, &r.UID, &r.Type, &r.RefID, &r.Provider, &r.StartDate, &r.EndDate, &r.CreatedAt,
			&r.CreatedBy, &r.RefUID, &r.ProviderUID); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// GetWarrantyByUID fetches a warranty by its UID
func GetWarrantyByUID(pool *pgxpool.Pool, uid string) (WarrantyRow, error) {
	ctx := context.Background()
	var r WarrantyRow
	err := pool.QueryRow(ctx, `
		SELECT id, uid, type, ref_id, provider, start_date, end_date, created_at, created_by
		FROM warranties
		WHERE uid = $1`, uid).
		Scan(&r.ID, &r.UID, &r.Type, &r.RefID, &r.Provider, &r.StartDate, &r.EndDate, &r.CreatedAt, &r.CreatedBy)
	return r, err
}

// CreateWarranty inserts a new warranty row
func CreateWarranty(pool *pgxpool.Pool, warrantyType string, refID int, provider *int, startDate, endDate string) (string, error) {
	ctx := context.Background()
	var uid string

	var startDateVal interface{} = nil
	if strings.TrimSpace(startDate) != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			startDateVal = t
		}
	}

	var endDateVal interface{} = nil
	if strings.TrimSpace(endDate) != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateVal = t
		}
	}
	if err := pool.QueryRow(ctx, `
		INSERT INTO warranties (type, ref_id, provider, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING uid`, warrantyType, refID, provider, startDateVal, endDateVal).Scan(&uid); err != nil {
		return "", err
	}
	return uid, nil
}

// DeleteWarranty removes a warranty by its ID
func DeleteWarranty(pool *pgxpool.Pool, id int) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM warranties WHERE id = $1`, id)
	return err
}

// UpdateWarranty applies partial updates to a warranty by ID
func UpdateWarranty(pool *pgxpool.Pool, id int, req WarrantyPatchReq) error {
	ctx := context.Background()
	setClauses := []string{}
	var args []any
	idx := 1

	if req.Provider != nil {
		setClauses = append(setClauses, fmt.Sprintf("provider = $%d", idx))
		args = append(args, *req.Provider)
		idx++
	}
	if req.StartDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("start_date = $%d", idx))
		args = append(args, *req.StartDate)
		idx++
	}
	if req.EndDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("end_date = $%d", idx))
		args = append(args, *req.EndDate)
		idx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE warranties SET %s WHERE id = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}
