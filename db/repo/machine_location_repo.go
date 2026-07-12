package repo

import (
	"context"
	"fmt"
	"strings"

	"app/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MachineLocationRow represents a machine_locations row
type MachineLocationRow struct {
	ID         int
	UID        pgtype.Text
	MachineID  int
	Plant      pgtype.Text
	Area       pgtype.Text
	CellOrLine pgtype.Text
	CreatedAt  pgtype.Timestamp
	CreatedBy  pgtype.Int8
}

// ListMachineLocations fetches locations for a machine
func ListMachineLocations(pool *pgxpool.Pool, machineID int) ([]MachineLocationRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `
		SELECT id, uid, machine_id, plant, area, cell_or_line, created_at, created_by
		FROM machine_locations
		WHERE machine_id = $1
		ORDER BY id ASC`, machineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []MachineLocationRow
	for rows.Next() {
		var r MachineLocationRow
		if err := rows.Scan(&r.ID, &r.UID, &r.MachineID, &r.Plant, &r.Area, &r.CellOrLine, &r.CreatedAt, &r.CreatedBy); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMachineLocationByUID fetches a location by its UID
func GetMachineLocationByUID(pool *pgxpool.Pool, uid string) (MachineLocationRow, error) {
	ctx := context.Background()
	var r MachineLocationRow
	err := pool.QueryRow(ctx, `
		SELECT id, uid, machine_id, plant, area, cell_or_line, created_at, created_by
		FROM machine_locations
		WHERE uid = $1`, uid).
		Scan(&r.ID, &r.UID, &r.MachineID, &r.Plant, &r.Area, &r.CellOrLine, &r.CreatedAt, &r.CreatedBy)
	return r, err
}

// CreateMachineLocation inserts a new machine location row
func CreateMachineLocation(pool *pgxpool.Pool, machineID int, plant, area, cellOrLine string) (string, error) {
	ctx := context.Background()
	var uid string
	if err := pool.QueryRow(ctx, `
		INSERT INTO machine_locations (machine_id, plant, area, cell_or_line)
		VALUES ($1, NULLIF($2, ''), NULLIF($3, ''), NULLIF($4, ''))
		RETURNING uid`, machineID, plant, area, cellOrLine).Scan(&uid); err != nil {
		return "", err
	}
	return uid, nil
}

// DeleteMachineLocation removes a location by its ID
func DeleteMachineLocation(pool *pgxpool.Pool, id int) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM machine_locations WHERE id = $1`, id)
	return err
}

// UpdateMachineLocation applies partial updates to a machine location by ID
func UpdateMachineLocation(pool *pgxpool.Pool, id int, req models.MachineLocationUpdateRequest) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	if req.Plant != nil {
		setClauses = append(setClauses, fmt.Sprintf("plant = $%d", idx))
		args = append(args, *req.Plant)
		idx++
	}
	if req.Area != nil {
		setClauses = append(setClauses, fmt.Sprintf("area = NULLIF($%d, '')", idx))
		args = append(args, *req.Area)
		idx++
	}
	if req.CellOrLine != nil {
		setClauses = append(setClauses, fmt.Sprintf("cell_or_line = NULLIF($%d, '')", idx))
		args = append(args, *req.CellOrLine)
		idx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE machine_locations SET %s WHERE id = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}
