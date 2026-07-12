package repo

import (
	models "app/models/maintenance"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TemplateProcedureRow struct {
	ID          int
	UID         string
	Name        string
	Description string
	TemplateID  string
	TemplateUID string
}

type TemplateProcedureQP struct {
	TemplateID   *int
	ProcedureUID *string
}

type TemplateProcedureInsert struct {
	Name        string
	Description string
	TemplateID  int
}

func GetTemplateProcedure(pool *pgxpool.Pool, qp TemplateProcedureQP) ([]TemplateProcedureRow, error) {
	ctx := context.Background()
	var filters []string
	var filterValues []any
	idx := 1
	query := `
		SELECT 
			mtp.id, 
			mtp.uid, 
			mtp.name, 
			mtp.description,
			mtp.maintenance_template_id,
			mt.uid maintenance_template_uid
		FROM 
			maintenance_template_procedures mtp
			join maintenance_templates mt on mt.id = mtp.maintenance_template_id 
	`

	if qp.TemplateID != nil {
		filters = append(filters, fmt.Sprintf(" mtp.maintenance_template_id = $%d ", idx))
		filterValues = append(filterValues, &qp.TemplateID)
		idx++
	}

	if qp.ProcedureUID != nil && *qp.ProcedureUID != "" {
		filters = append(filters, fmt.Sprintf(" mtp.uid = $%d ", idx))
		filterValues = append(filterValues, &qp.ProcedureUID)
		idx++
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += " ORDER BY mtp.id ASC"

	rows, err := pool.Query(ctx, query, filterValues...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var list []TemplateProcedureRow

	for rows.Next() {
		var r TemplateProcedureRow
		err := rows.Scan(&r.ID, &r.UID, &r.Name, &r.Description, &r.TemplateID, &r.TemplateUID)
		if err != nil {
			return nil, err
		}
		list = append(list, r)
	}

	return list, nil
}

func CreateTemplateProcedure(pool *pgxpool.Pool, d TemplateProcedureInsert) (string, error) {
	ctx := context.Background()
	var newUID string
	err := pool.QueryRow(ctx,
		`INSERT INTO maintenance_template_procedures (name, description, maintenance_template_id) VALUES ($1, $2, $3) RETURNING uid`,
		d.Name,
		d.Description,
		d.TemplateID,
	).Scan(&newUID)
	if err != nil {
		return "", err
	}
	return newUID, nil

}

func UpdateTemplateProcedure(pool *pgxpool.Pool, uid string, req models.TemplateProcedureUpdateRequest) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	// ---- NAME ----
	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, *req.Name)
		idx++
	}

	// ---- DESCRIPTION ----
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", idx))
		args = append(args, *req.Description)
		idx++
	}

	if req.TemplateUID != nil {
		setClauses = append(setClauses, fmt.Sprintf("maintenance_template_id = $%d", idx))
		args = append(args, *req.TemplateUID)
		idx++
	}

	// If nothing to update
	if len(setClauses) == 0 {
		return nil
	}

	// WHERE clause
	args = append(args, uid)
	query := "UPDATE maintenance_template_procedures SET " +
		strings.Join(setClauses, ", ") +
		" WHERE uid = " + fmt.Sprintf("$%d", idx)

	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeleteTemplateProcedure(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM maintenance_template_procedures WHERE uid = $1`, uid)
	return err
}
