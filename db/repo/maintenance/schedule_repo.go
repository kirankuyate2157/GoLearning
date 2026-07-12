package repo

import (
	models "app/models/maintenance"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MaintenanceScheduleRow struct {
	ID                     int
	UID                    string
	Date                   time.Time
	MaintenanceTemplateID  int
	MaintenanceTemplateUID string
	AssignedTo             *int
	MachineID              int
	MachineUID             string
	StartTimestamp         *time.Time
	EndTimestamp           *time.Time
	Findings               string
	Status                 string
	StartedOn              *time.Time
	EndedOn                *time.Time
}

type MaintenanceScheduleUpdate struct {
	Date                  *time.Time
	MaintenanceTemplateID *int
	AssignedTo            *int
	MachineID             *int
	StartTimestamp        *time.Time
	EndTimestamp          *time.Time
	Findings              *string
	Status                *string
	StartedOn             *time.Time
	EndedOn               *time.Time
}

func GetMaintenanceSchedule(
	pool *pgxpool.Pool,
	qp models.MaintenanceScheduleParams) ([]MaintenanceScheduleRow, error) {
	ctx := context.Background()
	var filters []string
	var filterValues []any
	idx := 1

	limit := qp.Limit
	offset := qp.Offset

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	query := `
        SELECT 
            ms.id,
            ms.uid,
            ms.maintenance_template_id,
            mt.uid AS maintenance_template_uid,
            ms.assigned_to,
            ms.machine_id,
            m.uid AS machine_uid,
            ms.start_timestamp,
            ms.end_timestamp,
            ms.findings,
            ms.status,
            ms.started_on,
            ms.ended_on
        FROM maintenance_schedules ms
        LEFT JOIN maintenance_templates mt ON ms.maintenance_template_id = mt.id
        LEFT JOIN machines m ON ms.machine_id = m.id
    `

	// Filters
	if qp.UID != "" {
		filters = append(filters, fmt.Sprintf("ms.uid = $%d", idx))
		filterValues = append(filterValues, qp.UID)
		idx++
	}

	if qp.Status != "" {
		filters = append(filters, fmt.Sprintf("ms.status = $%d", idx))
		filterValues = append(filterValues, qp.Status)
		idx++
	}

	// Apply filters
	whereClause := ""
	if len(filters) > 0 {
		whereClause = " WHERE " + strings.Join(filters, " AND ")
	}

	var err error

	query += whereClause
	query += " ORDER BY ms.date DESC, ms.id ASC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", idx, idx+1)

	filterValues = append(filterValues, limit, offset)

	rows, err := pool.Query(ctx, query, filterValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []MaintenanceScheduleRow
	for rows.Next() {
		var r MaintenanceScheduleRow
		var assignedTo, machineID, templateID sql.NullInt32
		var machineUID, templateUID, findings, status sql.NullString
		var startTimestamp, endTimestamp, startedOn, endedOn sql.NullTime

		if err := rows.Scan(
			&r.ID,
			&r.UID,
			&templateID,
			&templateUID,
			&assignedTo,
			&machineID,
			&machineUID,
			&startTimestamp,
			&endTimestamp,
			&findings,
			&status,
			&startedOn,
			&endedOn,
		); err != nil {
			return nil, err
		}

		if templateID.Valid {
			r.MaintenanceTemplateID = int(templateID.Int32)
		}
		if templateUID.Valid {
			r.MaintenanceTemplateUID = templateUID.String
		}
		if assignedTo.Valid {
			val := int(assignedTo.Int32)
			r.AssignedTo = &val
		}
		if machineID.Valid {
			r.MachineID = int(machineID.Int32)
		}
		if machineUID.Valid {
			r.MachineUID = machineUID.String
		}
		if startTimestamp.Valid {
			r.StartTimestamp = &startTimestamp.Time
		}
		if endTimestamp.Valid {
			r.EndTimestamp = &endTimestamp.Time
		}
		if findings.Valid {
			r.Findings = findings.String
		}
		if status.Valid {
			r.Status = status.String
		}
		if startedOn.Valid {
			r.StartedOn = &startedOn.Time
		}
		if endedOn.Valid {
			r.EndedOn = &endedOn.Time
		}

		out = append(out, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func CreateMaintenanceSchedule(
	pool *pgxpool.Pool,
	maintenanceTemplateID int,
	machineID int,
	startTimestamp,
	endTimestamp *time.Time) (string, error) {
	ctx := context.Background()
	var newUID string

	query := `
		INSERT INTO 
		maintenance_schedules (
			maintenance_template_id, 
			machine_id, 
			start_timestamp, 
			end_timestamp
		) 
		VALUES ($1, $2, $3, $4) RETURNING uid
	`
	err := pool.QueryRow(
		ctx,
		query,
		maintenanceTemplateID,
		machineID,
		startTimestamp,
		endTimestamp).Scan(&newUID)
	if err != nil {
		return "", err
	}
	return newUID, nil
}

func UpdateMaintenanceSchedule(pool *pgxpool.Pool, uid string, req MaintenanceScheduleUpdate) error {
	ctx := context.Background()

	var setClauses []string
	var args []any
	idx := 1

	if req.Date != nil {
		setClauses = append(setClauses, fmt.Sprintf("date = $%d", idx))
		args = append(args, *req.Date)
		idx++
	}

	if req.MaintenanceTemplateID != nil {
		setClauses = append(setClauses, fmt.Sprintf("maintenance_template_id = $%d", idx))
		args = append(args, *req.MaintenanceTemplateID)
		idx++
	}

	if req.AssignedTo != nil {
		setClauses = append(setClauses, fmt.Sprintf("assigned_to = $%d", idx))
		args = append(args, *req.AssignedTo)
		idx++
	}

	if req.MachineID != nil {
		setClauses = append(setClauses, fmt.Sprintf("machine_id = $%d", idx))
		args = append(args, *req.MachineID)
		idx++
	}

	if req.StartTimestamp != nil {
		setClauses = append(setClauses, fmt.Sprintf("start_timestamp = $%d", idx))
		args = append(args, *req.StartTimestamp)
		idx++
	}

	if req.EndTimestamp != nil {
		setClauses = append(setClauses, fmt.Sprintf("end_timestamp = $%d", idx))
		args = append(args, *req.EndTimestamp)
		idx++
	}

	if req.Findings != nil {
		setClauses = append(setClauses, fmt.Sprintf("findings = $%d", idx))
		args = append(args, *req.Findings)
		idx++
	}

	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", idx))
		args = append(args, *req.Status)
		idx++
	}

	if req.StartedOn != nil {
		setClauses = append(setClauses, fmt.Sprintf("started_on = $%d", idx))
		args = append(args, *req.StartedOn)
		idx++
	}

	if req.EndedOn != nil {
		setClauses = append(setClauses, fmt.Sprintf("ended_on = $%d", idx))
		args = append(args, *req.EndedOn)
		idx++
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := "UPDATE maintenance_schedules SET " + strings.Join(setClauses, ", ") +
		fmt.Sprintf(" WHERE uid = $%d", idx)
	args = append(args, uid)

	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeleteMaintenanceSchedule(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM maintenance_schedules WHERE uid = $1`, uid)
	return err
}
