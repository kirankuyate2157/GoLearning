package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaxesRow struct {
	ID        int
	UID       pgtype.Text
	Type      pgtype.Text
	Percent   float64
	CreatedAt pgtype.Timestamp
	CreatedBy pgtype.Text
}

type TaxesFilters struct {
	UID *string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type TaxesCreate struct {
	Type       string  `json:"type"`
	Percentage float64 `json:"percentage"`
}

type TaxesUpdate struct {
	Type       *string  `json:"type"`
	Percentage *float64 `json:"percentage"`
}

func ListTaxes(pool *pgxpool.Pool) ([]TaxesRow, error) {
	ctx := context.Background()

	query := `
		SELECT id,
		       uid,
		       type,
		       percent,
		       created_at,
		       created_by
		FROM tax
		ORDER BY id ASC
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []TaxesRow
	for rows.Next() {
		var r TaxesRow
		if err := rows.Scan(
			&r.ID,
			&r.UID,
			&r.Type,
			&r.Percent,
			&r.CreatedAt,
			&r.CreatedBy,
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

func GetTaxByUID(pool *pgxpool.Pool, uid string) (TaxesRow, error) {
	ctx := context.Background()
	var r TaxesRow
	err := pool.QueryRow(ctx, `SELECT id, uid, type, percent, created_at, created_by FROM tax WHERE uid = $1`, uid).Scan(&r.ID, &r.UID, &r.Type, &r.Percent, &r.CreatedAt, &r.CreatedBy)
	if err != nil {
		return r, err
	}
	return r, nil
}

func CreateTaxes(pool *pgxpool.Pool, req TaxesCreate) (string, error) {
	ctx := context.Background()
	var newUID string

	query := `
			INSERT INTO tax
			(type,
			 percent
			)
			VALUES ($1, $2) 
			RETURNING uid ;
	`
	err := pool.QueryRow(ctx, query,
		req.Type,
		req.Percentage).Scan(&newUID)
	if err != nil {
		return "", err
	}
	return newUID, nil
}

func UpdateTaxes(pool *pgxpool.Pool, id int, req TaxesUpdate) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	if req.Type != nil {
		setClauses = append(setClauses, fmt.Sprintf("type = $%d", idx))
		args = append(args, *req.Type)
		idx++
	}
	if req.Percentage != nil {
		setClauses = append(setClauses, fmt.Sprintf("percent = $%d", idx))
		args = append(args, *req.Percentage)
		idx++
	}
	if len(setClauses) == 0 {
		return nil
	}
	args = append(args, id)
	query := fmt.Sprintf("UPDATE tax SET %s WHERE id = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeleteTaxes(pool *pgxpool.Pool, id int) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM tax WHERE id = $1`, id)
	return err
}
