package repo

import (
	"app/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FeedTypeRow struct {
	ID          int
	UID         pgtype.Text
	Name        pgtype.Text
	Description pgtype.Text
	Strucuture  pgtype.Text
	CreatedAt   pgtype.Timestamp
}

// ListFeedTypes returns all feed connection types
func ListFeedTypes(pool *pgxpool.Pool) ([]FeedTypeRow, error) {
	ctx := context.Background()
	query := `SELECT id, 
					uid, 
					name, 
					description, 
					structure 
					FROM feed_connection_types 
					ORDER BY id ASC`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []FeedTypeRow
	for rows.Next() {
		var r FeedTypeRow
		if err := rows.Scan(&r.ID, &r.UID, &r.Name, &r.Description, &r.Strucuture); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFeedTypeByUID returns a single feed connection type by uid
func GetFeedTypeByUID(pool *pgxpool.Pool, uid string) (FeedTypeRow, error) {
	ctx := context.Background()
	query := `SELECT id, uid, name, description, structure, created_at FROM feed_connection_types WHERE uid = $1`
	var r FeedTypeRow
	err := pool.QueryRow(ctx, query, uid).Scan(&r.ID, &r.UID, &r.Name, &r.Description, &r.Strucuture, &r.CreatedAt)
	return r, err
}

// CreateFeedType inserts a new feed connection type and returns its UID
func CreateFeedType(pool *pgxpool.Pool, req models.FeedTypeCreateRequest) (string, error) {
	ctx := context.Background()

	// Convert structure map to JSON
	var structureJSON []byte
	var err error

	if req.Structure != nil {
		structureJSON, err = json.Marshal(req.Structure)
		if err != nil {
			return "", fmt.Errorf("failed to marshal structure: %w", err)
		}
	}

	var newUID string

	err = pool.QueryRow(ctx,
		`INSERT INTO feed_connection_types (name, description, structure)
		 VALUES ($1, $2, $3)
		 RETURNING uid`,
		req.Name,
		req.Description,
		structureJSON,
	).Scan(&newUID)

	if err != nil {
		return "", err
	}

	return newUID, nil
}

// UpdateFeedType applies partial updates to a feed connection type
func UpdateFeedType(pool *pgxpool.Pool, uid string, req models.FeedTypeUpdateRequest) error {
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

	// structure (JSON)
	if req.Structure != nil {
		var structureJSON []byte
		var err error

		if req.Structure != nil {
			structureJSON, err = json.Marshal(req.Structure)
			if err != nil {
				return fmt.Errorf("invalid structure JSON: %w", err)
			}
		}

		setClauses = append(setClauses, fmt.Sprintf("structure = $%d", idx))
		args = append(args, nullIfNilJSON(structureJSON))
		idx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, uid)

	query := fmt.Sprintf(`
		UPDATE feed_connection_types
		SET %s
		WHERE uid = $%d
	`, strings.Join(setClauses, ", "), idx)

	_, err := pool.Exec(ctx, query, args...)
	return err
}

// DeleteFeedType removes a feed connection type by UID
func DeleteFeedType(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM feed_connection_types WHERE uid = $1`, uid)
	return err
}

// NOTE: nullIfEmpty is defined in other repo files; reuse existing helper.
// nullIfEmpty returns nil if empty string.
func nullIfEmpty(s string) interface{} {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return s
}

// nullIfNilJSON returns nil if JSON byte array is nil.
func nullIfNilJSON(b []byte) interface{} {
	if b == nil {
		return nil
	}
	return b
}
