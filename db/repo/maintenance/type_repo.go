package repo

import (
	models "app/models/maintenance"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MaintenanceTypeRow struct {
	ID          int
	UID         pgtype.Text
	Name        pgtype.Text
	Description pgtype.Text
	CreatedAt   pgtype.Timestamp
}

func GetMaintenanceType(pool *pgxpool.Pool, qp models.MaintenanceType) ([]MaintenanceTypeRow, error) {
	ctx := context.Background()
	var filters []string
	var filterValues []any
	idx := 1
	query := `SELECT id, 
					uid, 
					name, 
					description 
					FROM maintenance_types mt`

	if len(qp.UID) != 0 {
		filters = append(filters, fmt.Sprintf(" mt.uid = $%d ", idx))
		filterValues = append(filterValues, qp.UID)
		idx++
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += " ORDER BY mt.id ASC"

	rows, err := pool.Query(ctx, query, filterValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []MaintenanceTypeRow
	for rows.Next() {
		var r MaintenanceTypeRow
		if err := rows.Scan(&r.ID, &r.UID, &r.Name, &r.Description); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateMaintenanceType inserts a new Maintenance connection type and returns its UID
func CreateMaintenanceType(pool *pgxpool.Pool, req models.MaintenanceTypeCreateRequest) (string, error) {
	ctx := context.Background()

	// Convert structure map to JSON
	var err error

	var newUID string

	err = pool.QueryRow(ctx,
		`INSERT INTO maintenance_types (name, description)
		 VALUES ($1, $2)
		 RETURNING uid`,
		req.Name,
		req.Description,
	).Scan(&newUID)

	if err != nil {
		return "", err
	}

	return newUID, nil
}

// UpdateMaintenanceType applies partial updates to a Maintenance connection type
func UpdateMaintenanceType(pool *pgxpool.Pool, uid string, req models.MaintenanceTypeUpdateRequest) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	// name
	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Name))
		idx++
	}

	// description
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", idx))
		args = append(args, nullIfEmpty(*req.Description))
		idx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, uid)

	query := fmt.Sprintf(`
		UPDATE maintenance_types
		SET %s
		WHERE uid = $%d
	`, strings.Join(setClauses, ", "), idx)

	_, err := pool.Exec(ctx, query, args...)
	return err
}

// DeleteMaintenanceType removes a Maintenance connection type by UID
func DeleteMaintenanceType(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM maintenance_types WHERE uid = $1`, uid)
	return err
}

func nullIfEmpty(s string) interface{} {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return s
}
