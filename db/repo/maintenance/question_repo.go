package repo

import (
	"context"
	"fmt"
	"strings"

	// "time"

	models "app/models/maintenance"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// StepQuestionRow represents a step question row from the database
type StepQuestionRow struct {
	ID           int
	UID          pgtype.Text
	StepUID      pgtype.Text
	Name         pgtype.Text
	Description  pgtype.Text
	Sequence     int
	DisplayOrder int
}

func CreateProcedureStepQuestion(pool *pgxpool.Pool,
	name string,
	description string,
	stepID int,
	sequence int,
	displayOrder int,
) (string, error) {
	ctx := context.Background()
	var uid string

	if err := pool.QueryRow(ctx, `
		INSERT INTO maintenance_procedure_steps_questions (name, description, step_id, sequence, display_order)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING uid`, name, description, stepID, sequence, displayOrder).Scan(&uid); err != nil {
		return "", err
	}
	return uid, nil
}

func GetProcedureStepQuestions(pool *pgxpool.Pool, qp models.ProcedureStepQuestionParams) ([]StepQuestionRow, error) {
	ctx := context.Background()
	var filters []string
	var filterValues []any
	idx := 1
	query := `
		SELECT 
			mpq.id,
			mpq.uid,
			mpq.name,
			mpq.description,
			mps.uid as procedure_step_uid,
			mpq.sequence,
			mpq.display_order
		FROM maintenance_procedure_steps_questions mpq
		LEFT JOIN maintenance_procedure_steps mps ON mpq.step_id = mps.id
	`

	if len(qp.UID) != 0 {
		filters = append(filters, fmt.Sprintf(" mpq.uid = $%d ", idx))
		filterValues = append(filterValues, qp.UID)
		idx++
	}

	if len(qp.StepUID) != 0 {
		filters = append(filters, fmt.Sprintf(" mps.uid = $%d ", idx))
		filterValues = append(filterValues, qp.StepUID)
		idx++
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += " ORDER BY mpq.display_order NULLS LAST, mpq.sequence, mpq.id"

	rows, err := pool.Query(ctx, query, filterValues...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var out []StepQuestionRow
	for rows.Next() {
		var r StepQuestionRow
		if err := rows.Scan(&r.ID,
			&r.UID,
			&r.Name,
			&r.Description,
			&r.StepUID,
			&r.Sequence,
			&r.DisplayOrder); err != nil {
			return nil, err
		}
		out = append(out, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil

}

// DeleteProcedureStepQuestion removes step question by its ID
func DeleteProcedureStepQuestion(pool *pgxpool.Pool, id int) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM maintenance_procedure_steps_questions WHERE id = $1`, id)
	return err
}

func UpdateProcedureStepQuestion(pool *pgxpool.Pool, id int, req models.ProcedureStepQuestionUpdateRequest) error {
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
	if req.Sequence != nil {
		setClauses = append(setClauses, fmt.Sprintf("sequence = $%d", idx))
		args = append(args, *req.Sequence)
		idx++
	}
	if req.DisplayOrder != nil {
		setClauses = append(setClauses, fmt.Sprintf("display_order = $%d", idx))
		args = append(args, *req.DisplayOrder)
		idx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf(
		"UPDATE maintenance_procedure_steps_questions SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "),
		idx,
	)
	_, err := pool.Exec(ctx, query, args...)
	return err
}
