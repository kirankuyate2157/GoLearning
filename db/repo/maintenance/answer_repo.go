package repo

import (
	models "app/models/maintenance"
	"context"
	"fmt"
	"strings"

	// models "app/models/maintenance"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// QuestionAnswerRow represents a question answer row from the database
type QuestionAnswerRow struct {
	ID           int
	UID          pgtype.Text
	QuestionUID  pgtype.Text
	AnswerText   pgtype.Text
	DisplayOrder int
}

type QustionAnswerFilters struct {
	UID         *string
	QuestionID  *int
	QuestionUID *string
}

func CreateQuestionAnswer(pool *pgxpool.Pool,
	questionID int,
	answerText string,
	displayOrder int,
) (string, error) {
	ctx := context.Background()
	var uid string

	if err := pool.QueryRow(ctx, `
		INSERT INTO maintenance_question_answers (question_id, answer_text, display_order)
		VALUES ($1, $2, $3)
		RETURNING uid`, questionID, answerText, displayOrder).Scan(&uid); err != nil {
		return "", err
	}
	return uid, nil
}

func GetQuestionAnswers(pool *pgxpool.Pool, qp QustionAnswerFilters) ([]QuestionAnswerRow, error) {
	ctx := context.Background()
	var filters []string
	var filterValues []any
	idx := 1
	query := `
		SELECT 
			qa.id,
			qa.uid,
			qa.answer_text,
			q.uid as question_uid,
			qa.display_order
		FROM maintenance_question_answers qa
		LEFT JOIN maintenance_procedure_steps_questions q ON qa.question_id = q.id
	`

	if qp.UID != nil && len(*qp.UID) != 0 {
		filters = append(filters, fmt.Sprintf(" qa.uid = $%d ", idx))
		filterValues = append(filterValues, *qp.UID)
		idx++
	}

	if qp.QuestionID != nil {
		filters = append(filters, fmt.Sprintf(" qa.question_id = $%d ", idx))
		filterValues = append(filterValues, *qp.QuestionID)
		idx++
	}

	if qp.QuestionUID != nil && len(*qp.QuestionUID) != 0 {
		filters = append(filters, fmt.Sprintf(" q.uid = $%d ", idx))
		filterValues = append(filterValues, *qp.QuestionUID)
		idx++
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += " ORDER BY qa.display_order NULLS LAST, qa.id"

	rows, err := pool.Query(ctx, query, filterValues...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var out []QuestionAnswerRow
	for rows.Next() {
		var r QuestionAnswerRow
		if err := rows.Scan(&r.ID,
			&r.UID,
			&r.AnswerText,
			&r.QuestionUID,
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

// DeleteQuestionAnswer removes question answer by its ID
func DeleteQuestionAnswer(pool *pgxpool.Pool, id int) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM maintenance_question_answers WHERE id = $1`, id)
	return err
}

func UpdateQuestionAnswer(pool *pgxpool.Pool, id int, req models.QuestionAnswerUpdateRequest) error {
	ctx := context.Background()
	setClauses := []string{}
	var args []any
	idx := 1

	if req.AnswerText != nil {
		setClauses = append(setClauses, fmt.Sprintf("answer_text = $%d", idx))
		args = append(args, *req.AnswerText)
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
		"UPDATE maintenance_question_answers SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "),
		idx,
	)
	_, err := pool.Exec(ctx, query, args...)
	return err
}
