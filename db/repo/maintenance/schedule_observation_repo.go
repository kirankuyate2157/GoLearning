package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MaintenanceScheduleObservationRow struct {
	ID                     int
	UID                    pgtype.Text
	MaintenanceScheduleUID pgtype.Text
	QuestionUID            pgtype.Text
	Remark                 pgtype.Text
	IsComplete             bool
	AnswerUID              pgtype.Text
}

type MaintenanceScheduleObservationCreateRequest struct {
	Remark                *string
	MaintenanceScheduleID *int
	AnswerID              *int
	QuestionID            *int
}

type MaintenanceScheduleObservationParams struct {
	UID                    *string
	MaintenanceScheduleUID *string
	QuestionUID            *string
}

// GetMaintenaneSchedulesChecklist retrieves maintenance schedule checklists based on provided filters
func GetMaintenaneScheduleQuestions(pool *pgxpool.Pool, qp MaintenanceScheduleObservationParams) ([]MaintenanceScheduleObservationRow, error) {
	ctx := context.Background()
	var filters []string
	var filterValues []any
	idx := 1

	query := `
		SELECT 
			mc.id,
			mc.uid,
			ms.uid as maintenance_schedule_uid,
			mpq.uid as question_uid,
			mc.remark,
			mc.is_complete,
			mqa.uid as answer_uid
		FROM 
		maintenance_schedule_questions mc 
		JOIN maintenance_schedules ms ON mc.maintenance_schedule_id = ms.id
		JOIN maintenance_procedure_steps_questions mpq ON mc.question_id = mpq.id
		LEFT JOIN maintenance_question_answers mqa ON mc.answer_id = mqa.id
	`

	// if len(qp.ProcedureStepChecklistUID) != 0 {
	if qp.QuestionUID != nil {
		filters = append(filters, fmt.Sprintf(" mpq.uid = $%d ", idx))
		filterValues = append(filterValues, *qp.QuestionUID)
		idx++
	}

	if qp.MaintenanceScheduleUID != nil {
		filters = append(filters, fmt.Sprintf(" ms.uid = $%d ", idx))
		filterValues = append(filterValues, *qp.MaintenanceScheduleUID)
		idx++
	}

	if qp.UID != nil {
		filters = append(filters, fmt.Sprintf(" mc.uid = $%d ", idx))
		filterValues = append(filterValues, *qp.UID)
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

	var out []MaintenanceScheduleObservationRow
	for rows.Next() {
		var r MaintenanceScheduleObservationRow
		if err := rows.Scan(&r.ID,
			&r.UID,
			&r.MaintenanceScheduleUID,
			&r.QuestionUID,
			&r.Remark,
			&r.IsComplete,
			&r.AnswerUID,
		); err != nil {
			return nil, err
		}
		out = append(out, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

// CreateMaintenanceScheduleObservation creates a mapping of maintenance schedule and step question
func CreateMaintenanceScheduleObservation(pool *pgxpool.Pool, req MaintenanceScheduleObservationCreateRequest) (string, error) {
	ctx := context.Background()
	var uid string

	query := `
			INSERT INTO maintenance_schedule_questions
				(maintenance_schedule_id, question_id, remark, answer_id)
			VALUES
				($1,$2,$3,$4)
			RETURNING uid;
			
	`
	err := pool.QueryRow(ctx, query,
		req.MaintenanceScheduleID,
		req.QuestionID,
		req.Remark,
		req.AnswerID,
	).Scan(&uid)

	if err != nil {
		return "", fmt.Errorf("insert feed failed: %w", err)
	}

	return uid, nil

}

// DeleteMaintenanceScheduleObservation removes a mapping of maintenance schedule and procedure step question by its ID
func DeleteMaintenanceScheduleObservation(pool *pgxpool.Pool, questionID int, scheduleID int) error {
	ctx := context.Background()

	query := `
		DELETE FROM maintenance_schedule_questions
		WHERE question_id = $1 AND maintenance_schedule_id = $2
	`

	_, err := pool.Exec(ctx, query, questionID, scheduleID)
	return err
}
