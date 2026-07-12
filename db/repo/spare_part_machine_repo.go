package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SparePartMachineRow represents a row from spare_parts_machines table
type SparePartMachineRow struct {
	ID           int
	UID          pgtype.Text
	SparePartID  int
	SparePartUID pgtype.Text
	MachineID    int
	MachineUID   pgtype.Text
	CreatedAt    pgtype.Timestamp
	CreatedBy    pgtype.Int8
}

// ListSparePartMachines lists machines for a given spare part (ordered by id)
func ListSparePartMachines(pool *pgxpool.Pool, sparePartUID string) ([]SparePartMachineRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `
	SELECT 
 	spm.id,
  	spm.uid,
  	spm.spare_part_id,
  	s.uid AS spare_part_uid,
  	spm.machine_id,
	m.uid AS machine_uid,
  	spm.created_at,
 	spm.created_by
	FROM spare_parts_machines spm
	JOIN spare_parts s ON spm.spare_part_id = s.id
	JOIN machines m ON spm.machine_id = m.id
	WHERE s.uid = $1
	ORDER BY spm.id ASC
`, sparePartUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []SparePartMachineRow
	for rows.Next() {
		var r SparePartMachineRow
		if err := rows.Scan(
			&r.ID,
			&r.UID,
			&r.SparePartID,
			&r.SparePartUID,
			&r.MachineID,
			&r.MachineUID,
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

// CreateSparePartMachine inserts a new mapping; expects sparePartID and machineID (int) resolved by caller

func CreateSparePartMachine(pool *pgxpool.Pool, sparePartID int, machineID int) (string, error) {
	ctx := context.Background()
	var newUID string
	err := pool.QueryRow(ctx, `
		INSERT INTO spare_parts_machines (spare_part_id, machine_id)
		VALUES ($1, $2)
		RETURNING uid
`, sparePartID, machineID).Scan(&newUID)
	if err != nil {
		return "", err
	}
	return newUID, nil
}

func DeleteSparePartMachine(pool *pgxpool.Pool, sparePartID int, machineID int) (int64, error) {
	cmd, err := pool.Exec(context.Background(),
		`DELETE FROM spare_parts_machines 
         WHERE spare_part_id=$1 AND machine_id=$2`,
		sparePartID, machineID,
	)
	if err != nil {
		return 0, err
	}
	return cmd.RowsAffected(), nil
}
