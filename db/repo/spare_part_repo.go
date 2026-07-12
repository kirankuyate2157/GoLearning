package repo

import (
	"context"
	"fmt"
	"strings"

	"app/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SparePartRow represents a spare part row from the database
type SparePartRow struct {
	ID          int
	UID         pgtype.Text
	Name        string
	Code        string
	Description pgtype.Text
	CreatedAt   pgtype.Timestamp
	CreatedBy   pgtype.Int8
	TypeUID     pgtype.Text
}

// ListSpareParts returns all spare parts ordered by ID
func ListSpareParts(pool *pgxpool.Pool) ([]SparePartRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `
	SELECT 
    s.id,
    s.uid,
    s.name,
    s.code,
    s.description,
    s.created_at,
    s.created_by,
    t.uid AS type_uid
FROM spare_parts s
LEFT JOIN types t ON s.type_id = t.id
ORDER BY s.id ASC

		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []SparePartRow
	for rows.Next() {
		var r SparePartRow
		if err := rows.Scan(&r.ID, &r.UID, &r.Name, &r.Code, &r.Description, &r.CreatedAt, &r.CreatedBy, &r.TypeUID); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// GetSparePartByUID returns a single spare part by its UID
func GetSparePartByUID(pool *pgxpool.Pool, uid string) (SparePartRow, error) {
	ctx := context.Background()
	var r SparePartRow
	err := pool.QueryRow(ctx, `
		SELECT 
    s.id,
    s.uid,
    s.name,
    s.code,
    s.description,
    s.created_at,
    s.created_by,
    t.uid AS type_uid
FROM spare_parts s
LEFT JOIN types t ON s.type_id = t.id
WHERE s.uid = $1`, uid).
		Scan(&r.ID, &r.UID, &r.Name, &r.Code, &r.Description, &r.CreatedAt, &r.CreatedBy, &r.TypeUID)
	return r, err
}

// CreateSparePart inserts a new spare part and returns its UID
func CreateSparePart(pool *pgxpool.Pool, req models.SparePartCreateRequest, typeID *int) (string, error) {
	ctx := context.Background()
	var newUID string

	err := pool.QueryRow(ctx, `
		INSERT INTO spare_parts (name, code, description , type_id)
		VALUES ($1, $2,  NULLIF($3, ''),$4)
		RETURNING uid`,
		strings.TrimSpace(req.Name),
		strings.TrimSpace(req.Code),
		strings.TrimSpace(req.Description),
		typeID).
		Scan(&newUID)
	if err != nil {
		return "", err
	}
	return newUID, nil
}

// UpdateSparePart applies partial updates to a spare part identified by UID
func UpdateSparePart(pool *pgxpool.Pool, uid string, req models.SparePartUpdateRequest, typeID *int) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Name))
		idx++
	}
	if req.Code != nil {
		setClauses = append(setClauses, fmt.Sprintf("code = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Code))
		idx++
	}
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = NULLIF($%d, '')", idx))
		args = append(args, strings.TrimSpace(*req.Description))
		idx++
	}

	if typeID != nil {
		setClauses = append(setClauses, fmt.Sprintf("type_id = $%d", idx))
		args = append(args, *typeID)
		idx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, uid)
	query := fmt.Sprintf("UPDATE spare_parts SET %s WHERE uid = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

// DeleteSparePart removes a spare part by UID
func DeleteSparePart(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM spare_parts WHERE uid = $1`, uid)
	return err
}
