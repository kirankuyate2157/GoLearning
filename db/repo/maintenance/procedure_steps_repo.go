package repo

import (
	"context"
	"fmt"
	"strings"

	models "app/models/maintenance"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProcedureStepRow struct {
	ID                              int
	UID                             string
	Name                            string
	Description                     string
	StepNo                          int
	MaintenanceTemplateProcedureUID string
}

// List Procedure Steps
func GetProcedureSteps(pool *pgxpool.Pool, qp models.ProcedureStepParams) ([]ProcedureStepRow, error) {
	ctx := context.Background()
	var filters []string
	var filterValues []any
	idx := 1
	query := `
		SELECT 
			mps.id,
			mps.uid,
			mps.name,
			mps.description,
			mtp.uid as maintenance_template_procedure_uid,
			mps.step_no
		FROM maintenance_procedure_steps mps
		JOIN maintenance_template_procedures mtp
		ON mps.maintenance_template_procedure_id = mtp.id
	`
	if len(qp.ProcedureUID) != 0 {
		filters = append(filters, fmt.Sprintf(" mtp.uid = $%d ", idx))
		filterValues = append(filterValues, qp.ProcedureUID)
		idx++
	}

	if len(qp.UID) != 0 {
		filters = append(filters, fmt.Sprintf(" mps.uid = $%d ", idx))
		filterValues = append(filterValues, qp.UID)
		idx++
	}

	if qp.StepNo != 0 {
		filters = append(filters, fmt.Sprintf(" mps.step_no = $%d ", idx))
		filterValues = append(filterValues, qp.StepNo)
		idx++
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += " ORDER BY mps.id ASC"

	rows, err := pool.Query(ctx, query, filterValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ProcedureStepRow
	for rows.Next() {
		var r ProcedureStepRow
		if err := rows.Scan(&r.ID, &r.UID, &r.Name, &r.Description, &r.MaintenanceTemplateProcedureUID, &r.StepNo); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// create Procdeure steps
func CreateProcedureSteps(pool *pgxpool.Pool, req models.ProcedureStepCreate, procedureID int) (string, error) {
	ctx := context.Background()
	var uid string

	query := `
		INSERT INTO maintenance_procedure_steps 
        	(name, description,maintenance_template_procedure_id,step_no)
        VALUES
        	($1,$2,$3,$4)
        RETURNING uid;
		`
	err := pool.QueryRow(ctx, query,
		req.Name,
		req.Description,
		procedureID,
		req.StepNo,
	).Scan(&uid)

	if err != nil {
		return "", fmt.Errorf("insert feed failed: %w", err)
	}

	return uid, nil
}

func UpdateProcedureStep(pool *pgxpool.Pool, stepID int, req models.ProcedureStepUpdate) error {
	ctx := context.Background()
	setClauses := []string{}
	var args []any
	idx := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, *req.Name)
		idx++
	}
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", idx))
		args = append(args, *req.Description)
		idx++
	}
	if req.StepNo != nil {
		setClauses = append(setClauses, fmt.Sprintf("step_no = $%d", idx))
		args = append(args, *req.StepNo)
		idx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, stepID)
	query := fmt.Sprintf(
		"UPDATE maintenance_procedure_steps SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "),
		idx,
	)
	_, err := pool.Exec(ctx, query, args...)

	return err
}

// DeleteProcedureSteps removes a step by its ID
func DeleteProcedureSteps(pool *pgxpool.Pool, id int) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM maintenance_procedure_steps WHERE id = $1`, id)
	return err
}
