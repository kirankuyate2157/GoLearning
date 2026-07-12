package repo

import (
	models "app/models/maintenance"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MaintenanceTemplateRow struct {
	ID          int
	UID         string
	Name        string
	TypeID      *int
	TypeUID     *string `json:"type_uid"`
	Description string
	Frequency   string
}

type MaintenanceTemplateUpdate struct {
	Name        *string
	Description *string
	FrequencyID *int
	TypeID      *int
}

func GetMaintenanceTemplate(pool *pgxpool.Pool, qp models.MaintenanceTemplateParams, withTotal bool) ([]MaintenanceTemplateRow, int, error) {
	ctx := context.Background()
	var filters []string
	var filterValues []any
	idx := 1

	limit := qp.Limit
	offset := qp.Offset

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	query := `
        SELECT 
            mtemp.id,
            mtemp.uid,
            mtemp.name,
            mt.uid AS type_uid,
            mtemp.description,
            t.uid as frequency
        FROM maintenance_templates mtemp
        JOIN maintenance_types mt
            ON mtemp.type_id = mt.id
		JOIN types t  
			ON mtemp.frequency = t.id
    `

	// Filters
	if qp.UID != "" {
		filters = append(filters, fmt.Sprintf("mtemp.uid = $%d", idx))
		filterValues = append(filterValues, qp.UID)
		idx++
	}

	if qp.ID != 0 {
		filters = append(filters, fmt.Sprintf("mtemp.id = $%d", idx))
		filterValues = append(filterValues, qp.ID)
		idx++
	}

	// Apply filters
	whereClause := ""
	if len(filters) > 0 {
		whereClause = " WHERE " + strings.Join(filters, " AND ")
	}

	var total int
	var err error
	if withTotal {
		countQuery := `SELECT COUNT(*) FROM maintenance_templates mtemp` + whereClause
		err = pool.QueryRow(ctx, countQuery, filterValues...).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	}

	query += whereClause
	query += " ORDER BY mtemp.id ASC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", idx, idx+1)

	filterValues = append(filterValues, limit, offset)

	rows, err := pool.Query(ctx, query, filterValues...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []MaintenanceTemplateRow
	for rows.Next() {
		var r MaintenanceTemplateRow
		var typeUID sql.NullString

		if err := rows.Scan(&r.ID, &r.UID, &r.Name, &typeUID, &r.Description, &r.Frequency); err != nil {
			return nil, 0, err
		}

		if typeUID.Valid {
			tid := typeUID.String
			r.TypeUID = &tid
		}

		out = append(out, r)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	if withTotal {
		return out, total, nil
	}

	return out, 0, nil
}

func CreateMaintenanceTemplate(pool *pgxpool.Pool, name, description string, type_, frequency int) (string, error) {
	ctx := context.Background()
	var newUID string

	err := pool.QueryRow(ctx,
		`INSERT INTO maintenance_templates (name, type_id, description, frequency) VALUES ($1, $2, $3, $4) RETURNING uid`,
		name,
		type_,
		description,
		frequency,
	).Scan(&newUID)
	if err != nil {
		return "", err
	}
	return newUID, nil
}

func UpdateMaintenanceTemplate(pool *pgxpool.Pool, uid string, req MaintenanceTemplateUpdate) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []any{}
	idx := 1

	// ---- NAME ----
	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, *req.Name)
		idx++
	}

	// ---- TYPE_UID ----
	if req.TypeID != nil {

		setClauses = append(setClauses, fmt.Sprintf("type_id = $%d", idx))
		args = append(args, req.TypeID)
		idx++

	}

	// ---- DESCRIPTION ----
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", idx))
		args = append(args, *req.Description)
		idx++
	}

	// ---- FREQUENCY ----
	if req.FrequencyID != nil {
		setClauses = append(setClauses, fmt.Sprintf("frequency = $%d", idx))
		args = append(args, *req.FrequencyID)
		idx++
	}

	// Nothing to update
	if len(setClauses) == 0 {
		return nil
	}

	// WHERE clause
	args = append(args, uid)
	query := "UPDATE maintenance_templates SET " + joinClauses(setClauses) +
		" WHERE uid = " + fmt.Sprintf("$%d", idx)

	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeleteMaintenanceTemplate(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM maintenance_templates WHERE uid = $1`, uid)
	return err
}

func joinClauses(clauses []string) string {
	return strings.Join(clauses, ", ")
}
