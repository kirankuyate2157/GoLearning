package repo

import (
	"context"
	// "encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlertTypeRow struct {
	ID          int
	UID         pgtype.Text
	Name        pgtype.Text
	Description pgtype.Text
	Structure   pgtype.Text
	CreatedAt   pgtype.Timestamp
}

type AlertsRow struct {
	ID        int
	UID       pgtype.Text
	Name      pgtype.Text
	Config    pgtype.Text
	TypeUID   pgtype.Text
	CreatedAt pgtype.Timestamp
}

type AlertsCreate struct {
	Name      string                   `json:"name"`
	TypeID    int                      `json:"type_id"`
	Structure []map[string]interface{} `json:"structure"`
}

type AlertsUpdate struct {
	Name      *string                   `json:"name"`
	TypeID    *int                      `json:"type_id"`
	Structure *[]map[string]interface{} `json:"structure"`
}

type AlertsTypeFilters struct {
	UID *string `json:"typeUID"  example:"550e8400-e29b-41d4-a716-446655440000"`
}

type AlertsFilters struct {
	UID *string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
}

func ListAlertsTypes(pool *pgxpool.Pool, qp AlertsTypeFilters) ([]AlertTypeRow, error) {
	ctx := context.Background()
	var filters []string
	var filterValues []any
	idx := 1

	query := `SELECT id, 
					uid, 
					name, 
					description, 
					structure 
					FROM alert_types at
					`

	// Filters
	if qp.UID != nil {
		filters = append(filters, fmt.Sprintf("at.uid = $%d", idx))
		filterValues = append(filterValues, qp.UID)
		idx++
	}

	// Apply filters
	whereClause := ""
	if len(filters) > 0 {
		whereClause = " WHERE " + strings.Join(filters, " AND ")
	}

	query += whereClause
	query += " ORDER BY at.id ASC"

	rows, err := pool.Query(ctx, query, filterValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []AlertTypeRow
	for rows.Next() {
		var r AlertTypeRow
		if err := rows.Scan(&r.ID, &r.UID, &r.Name, &r.Description, &r.Structure); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func ListAlerts(pool *pgxpool.Pool, qp AlertsFilters) ([]AlertsRow, error) {
	ctx := context.Background()
	var filters []string
	var filterValues []any
	idx := 1

	query := `SELECT a.id, 
					a.uid, 
					a.name, 
					a.config,
					at.uid as type_uid 
					FROM alerts a
					JOIN alert_types at ON a.type_id = at.id
					`
	// Filters
	if qp.UID != nil {
		filters = append(filters, fmt.Sprintf("a.uid = $%d", idx))
		filterValues = append(filterValues, qp.UID)
		idx++
	}

	// Apply filters
	whereClause := ""
	if len(filters) > 0 {
		whereClause = " WHERE " + strings.Join(filters, " AND ")
	}

	query += whereClause
	query += " ORDER BY a.id ASC"

	rows, err := pool.Query(ctx, query, filterValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []AlertsRow
	for rows.Next() {
		var r AlertsRow
		if err := rows.Scan(&r.ID,
			&r.UID,
			&r.Name,
			&r.Config,
			&r.TypeUID); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func CreateAlerts(pool *pgxpool.Pool, req AlertsCreate) (string, error) {
	ctx := context.Background()
	var newUID string

	query := `
			INSERT INTO alerts
			(name,
			 type_id,
			 config
			)
			VALUES ($1, $2, $3) 
			RETURNING uid ;
	`
	err := pool.QueryRow(ctx, query,
		req.Name,
		req.TypeID,
		req.Structure).Scan(&newUID)
	if err != nil {
		return "", err
	}
	return newUID, nil
}

func UpdateAlert(pool *pgxpool.Pool, id int, req AlertsUpdate) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, *req.Name)
		idx++
	}
	if req.TypeID != nil {
		setClauses = append(setClauses, fmt.Sprintf("type_id = $%d", idx))
		args = append(args, *req.TypeID)
		idx++
	}
	if req.Structure != nil {
		setClauses = append(setClauses, fmt.Sprintf("config = $%d", idx))
		args = append(args, *req.Structure)
		idx++
	}
	if len(setClauses) == 0 {
		return nil
	}
	args = append(args, id)
	query := fmt.Sprintf("UPDATE alerts SET %s WHERE id = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeleteAlert(pool *pgxpool.Pool, id int) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM alerts WHERE id = $1`, id)
	return err
}
