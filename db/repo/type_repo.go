package repo

import (
	"app/models"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TypeRow struct {
	ID          int
	UID         pgtype.Text
	Name        pgtype.Text
	Description pgtype.Text
	Category    pgtype.Text
	CreatedAt   pgtype.Timestamp
	CreatedBy   pgtype.Int8
}

func ListTypes(pool *pgxpool.Pool, category string) ([]TypeRow, error) {
	ctx := context.Background()
	query := `
		SELECT id, uid, name, description, category
		FROM types
	`
	var rows pgx.Rows
	var err error

	if category != "" {
		query += " WHERE category = $1 ORDER BY id ASC"
		rows, err = pool.Query(ctx, query, category)
	} else {
		query += " ORDER BY id ASC"
		rows, err = pool.Query(ctx, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []TypeRow
	for rows.Next() {
		var r TypeRow
		if err := rows.Scan(&r.ID, &r.UID, &r.Name, &r.Description, &r.Category); err != nil {
			return nil, err
		}
		out = append(out, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func ListTypeCategory(pool *pgxpool.Pool) ([]string, error) {
	ctx := context.Background()

	query := `
		SELECT e.enumlabel
		FROM pg_enum e
		JOIN pg_type t ON e.enumtypid = t.oid
		WHERE t.typname = 'types_category_enum'
		ORDER BY e.enumsortorder;
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		out = append(out, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

// GetTypeByUID returns a single machine by its UID
func GetTypeByUID(pool *pgxpool.Pool, uid string) (TypeRow, error) {
	ctx := context.Background()
	query := `
		SELECT 
			t.id,
			t.uid,
			t.name,
			t.description,
			t.category,
			t.created_at,
			t.created_by
		FROM types t
		WHERE t.uid = $1;
	`
	var r TypeRow
	err := pool.QueryRow(ctx, query, uid).Scan(
		&r.ID,
		&r.UID,
		&r.Name,
		&r.Description,
		&r.Category,
		&r.CreatedAt,
		&r.CreatedBy,
	)
	return r, err
}

// CreateType inserts a new type and returns its UID
func CreateType(pool *pgxpool.Pool, req models.TypeCreateRequest) (string, error) {
	ctx := context.Background()
	var newUID string
	err := pool.QueryRow(ctx,
		`INSERT INTO types (name,  description, category)
         VALUES ($1, $2, $3)
         RETURNING uid`,
		req.Name,
		nullIfEmpty(req.Description),
		req.Category,
	).Scan(&newUID)
	if err != nil {
		return "", err
	}
	return newUID, nil
}

// UpdateType applies partial updates to a type identified by UID
func UpdateType(pool *pgxpool.Pool, uid string, req models.TypeUpdateRequest) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Name))
		idx++
	}

	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", idx))
		args = append(args, nullIfEmpty(*req.Description))
		idx++
	}
	if req.Category != nil {
		setClauses = append(setClauses, fmt.Sprintf("category = $%d", idx))
		args = append(args, req.Category)
		idx++
	}
	// if req.CreatedBy != nil {
	// 	setClauses = append(setClauses, fmt.Sprintf("created_by = $%d", idx))
	// 	args = append(args, *req.CreatedBy)
	// 	idx++
	// }

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, uid)
	query := fmt.Sprintf("UPDATE types SET %s WHERE uid = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

// DeleteType removes a type by UID
func DeleteType(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM types WHERE uid = $1`, uid)
	return err
}
