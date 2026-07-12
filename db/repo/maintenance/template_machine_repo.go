package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MaintenanceTemplateMachineRow represents a row from maintenance_template_machines table
type MaintenanceTemplateMachineRow struct {
	ID                     int
	UID                    pgtype.Text
	MaintenanceTemplateID  int
	MaintenanceTemplateUID pgtype.Text
	MachineID              int
	MachineUID             pgtype.Text
	CreatedAt              pgtype.Timestamp
	CreatedBy              pgtype.Int8
}

// ListMaintenanceTemplateMachines lists machines for a given maintenance template (ordered by id)
func ListMaintenanceTemplateMachines(pool *pgxpool.Pool, maintenanceTemplateUID string) ([]MaintenanceTemplateMachineRow, error) {
	ctx := context.Background()
	query := `
		SELECT 
			mtm.id,
			mtm.uid,
			mtm.maintenance_template_id,
			mt.uid AS maintenance_template_uid,
			mtm.machine_id,
			m.uid AS machine_uid,
			mtm.created_at,
			mtm.created_by
		FROM maintenance_template_machines mtm
		JOIN maintenance_templates mt ON mtm.maintenance_template_id = mt.id
		JOIN machines m ON mtm.machine_id = m.id
		WHERE mt.uid = $1
		ORDER BY mtm.id ASC
	`
	rows, err := pool.Query(ctx, query, maintenanceTemplateUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []MaintenanceTemplateMachineRow
	for rows.Next() {
		var r MaintenanceTemplateMachineRow
		if err := rows.Scan(
			&r.ID,
			&r.UID,
			&r.MaintenanceTemplateID,
			&r.MaintenanceTemplateUID,
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

// CreateMaintenanceTemplateMachine inserts a new mapping; expects maintenanceTemplateID and machineID (int) resolved by caller
func CreateMaintenanceTemplateMachine(pool *pgxpool.Pool, maintenanceTemplateID int, machineID int) (string, error) {
	ctx := context.Background()
	var newUID string
	err := pool.QueryRow(ctx, `
		INSERT INTO maintenance_template_machines (maintenance_template_id, machine_id)
		VALUES ($1, $2)
		RETURNING uid
`, maintenanceTemplateID, machineID).Scan(&newUID)
	if err != nil {
		return "", err
	}
	return newUID, nil
}

// DeleteMaintenanceTemplateMachine deletes a mapping by maintenance template ID and machine ID
func DeleteMaintenanceTemplateMachine(pool *pgxpool.Pool, maintenanceTemplateID int, machineID int) (int64, error) {
	cmd, err := pool.Exec(context.Background(),
		`DELETE FROM maintenance_template_machines 
         WHERE maintenance_template_id=$1 AND machine_id=$2`,
		maintenanceTemplateID, machineID,
	)
	if err != nil {
		return 0, err
	}
	return cmd.RowsAffected(), nil
}

// GetMaintenanceTemplateMachineByUID retrieves a single mapping by UID
func GetMaintenanceTemplateMachineByUID(pool *pgxpool.Pool, uid string) (MaintenanceTemplateMachineRow, error) {
	ctx := context.Background()
	var r MaintenanceTemplateMachineRow
	err := pool.QueryRow(ctx, `
	SELECT 
 	mtm.id,
  	mtm.uid,
  	mtm.maintenance_template_id,
  	mt.uid AS maintenance_template_uid,
  	mtm.machine_id,
	m.uid AS machine_uid,
  	mtm.created_at,
 	mtm.created_by
	FROM maintenance_template_machines mtm
	JOIN maintenance_templates mt ON mtm.maintenance_template_id = mt.id
	JOIN machines m ON mtm.machine_id = m.id
	WHERE mtm.uid = $1
`, uid).Scan(
		&r.ID,
		&r.UID,
		&r.MaintenanceTemplateID,
		&r.MaintenanceTemplateUID,
		&r.MachineID,
		&r.MachineUID,
		&r.CreatedAt,
		&r.CreatedBy,
	)
	if err != nil {
		return MaintenanceTemplateMachineRow{}, err
	}
	return r, nil
}
