package repo

import (
	"app/models"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BOMRow struct {
	ID                int
	UID               pgtype.UUID
	ProductionOrderID pgtype.Int4
	JobID             pgtype.Int4
	Name              string
	PartNumber        pgtype.Text
	Description       pgtype.Text
	CreatedAt         pgtype.Timestamp
	UpdatedAt         pgtype.Timestamp
}

type BOMLineRow struct {
	ID         int
	UID        pgtype.UUID
	BOMID      int
	Name       string
	PartNumber pgtype.Text
	Quantity   float64
	Unit       pgtype.Text
	Notes      pgtype.Text
	CreatedAt  pgtype.Timestamp
	UpdatedAt  pgtype.Timestamp
}

type RoutingOperationRow struct {
	ID                     int
	UID                    pgtype.UUID
	ProductionOrderID      pgtype.Int4
	JobID                  pgtype.Int4
	SequenceNumber         int
	Name                   string
	Description            pgtype.Text
	EstimatedDurationHours pgtype.Numeric
	MachineID              pgtype.Int4
	Notes                  pgtype.Text
	CreatedAt              pgtype.Timestamp
	UpdatedAt              pgtype.Timestamp
}

func ListBOMs(pool *pgxpool.Pool) ([]BOMRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `SELECT id, uid, production_order_id, job_id, name, part_number, description, created_at, updated_at FROM boms ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []BOMRow
	for rows.Next() {
		var r BOMRow
		if err := rows.Scan(&r.ID, &r.UID, &r.ProductionOrderID, &r.JobID, &r.Name, &r.PartNumber, &r.Description, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func GetBOMByUID(pool *pgxpool.Pool, uid string) (BOMRow, error) {
	ctx := context.Background()
	var r BOMRow
	err := pool.QueryRow(ctx, `SELECT id, uid, production_order_id, job_id, name, part_number, description, created_at, updated_at FROM boms WHERE uid = $1`, uid).Scan(&r.ID, &r.UID, &r.ProductionOrderID, &r.JobID, &r.Name, &r.PartNumber, &r.Description, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}

func CreateBOM(pool *pgxpool.Pool, req models.BOMCreateRequest, productionOrderID *int, jobID *int) (string, error) {
	ctx := context.Background()
	var uid string
	err := pool.QueryRow(ctx, `INSERT INTO boms (production_order_id, job_id, name, part_number, description) VALUES ($1, $2, $3, $4, $5) RETURNING uid`, productionOrderID, jobID, strings.TrimSpace(req.Name), strings.TrimSpace(req.PartNumber), strings.TrimSpace(req.Description)).Scan(&uid)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func UpdateBOM(pool *pgxpool.Pool, uid string, req models.BOMUpdateRequest, productionOrderID *int, jobID *int) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1
	if req.ProductionOrderUID != nil {
		setClauses = append(setClauses, fmt.Sprintf("production_order_id = $%d", idx))
		args = append(args, productionOrderID)
		idx++
	}
	if req.JobUID != nil {
		setClauses = append(setClauses, fmt.Sprintf("job_id = $%d", idx))
		args = append(args, jobID)
		idx++
	}
	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Name))
		idx++
	}
	if req.PartNumber != nil {
		setClauses = append(setClauses, fmt.Sprintf("part_number = $%d", idx))
		args = append(args, strings.TrimSpace(*req.PartNumber))
		idx++
	}
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Description))
		idx++
	}
	if len(setClauses) == 0 {
		return nil
	}
	args = append(args, uid)
	query := fmt.Sprintf("UPDATE boms SET %s, updated_at = NOW() WHERE uid = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeleteBOM(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM boms WHERE uid = $1`, uid)
	return err
}

func ListBOMLines(pool *pgxpool.Pool, bomUID string) ([]BOMLineRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `SELECT bl.id, bl.uid, bl.bom_id, bl.name, bl.part_number, bl.quantity, bl.unit, bl.notes, bl.created_at, bl.updated_at FROM bom_lines bl JOIN boms b ON bl.bom_id = b.id WHERE b.uid = $1 ORDER BY bl.id ASC`, bomUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []BOMLineRow
	for rows.Next() {
		var r BOMLineRow
		if err := rows.Scan(&r.ID, &r.UID, &r.BOMID, &r.Name, &r.PartNumber, &r.Quantity, &r.Unit, &r.Notes, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func GetBOMLineByUID(pool *pgxpool.Pool, uid string) (BOMLineRow, error) {
	ctx := context.Background()
	var r BOMLineRow
	err := pool.QueryRow(ctx, `SELECT bl.id, bl.uid, bl.bom_id, bl.name, bl.part_number, bl.quantity, bl.unit, bl.notes, bl.created_at, bl.updated_at FROM bom_lines bl WHERE bl.uid = $1`, uid).Scan(&r.ID, &r.UID, &r.BOMID, &r.Name, &r.PartNumber, &r.Quantity, &r.Unit, &r.Notes, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}

func CreateBOMLine(pool *pgxpool.Pool, bomUID string, req models.BOMLineCreateRequest) (string, error) {
	ctx := context.Background()
	var uid string
	err := pool.QueryRow(ctx, `INSERT INTO bom_lines (bom_id, name, part_number, quantity, unit, notes) SELECT id, $2, $3, $4, $5, $6 FROM boms WHERE uid = $1 RETURNING uid`, bomUID, strings.TrimSpace(req.Name), strings.TrimSpace(req.PartNumber), req.Quantity, strings.TrimSpace(req.Unit), strings.TrimSpace(req.Notes)).Scan(&uid)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func UpdateBOMLine(pool *pgxpool.Pool, uid string, req models.BOMLineUpdateRequest) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1
	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Name))
		idx++
	}
	if req.PartNumber != nil {
		setClauses = append(setClauses, fmt.Sprintf("part_number = $%d", idx))
		args = append(args, strings.TrimSpace(*req.PartNumber))
		idx++
	}
	if req.Quantity != nil {
		setClauses = append(setClauses, fmt.Sprintf("quantity = $%d", idx))
		args = append(args, *req.Quantity)
		idx++
	}
	if req.Unit != nil {
		setClauses = append(setClauses, fmt.Sprintf("unit = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Unit))
		idx++
	}
	if req.Notes != nil {
		setClauses = append(setClauses, fmt.Sprintf("notes = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Notes))
		idx++
	}
	if len(setClauses) == 0 {
		return nil
	}
	args = append(args, uid)
	query := fmt.Sprintf("UPDATE bom_lines SET %s, updated_at = NOW() WHERE uid = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeleteBOMLine(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM bom_lines WHERE uid = $1`, uid)
	return err
}

func ListRoutingOperations(pool *pgxpool.Pool) ([]RoutingOperationRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `SELECT id, uid, production_order_id, job_id, sequence_number, name, description, estimated_duration_hours, machine_id, notes, created_at, updated_at FROM routing_operations ORDER BY sequence_number ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []RoutingOperationRow
	for rows.Next() {
		var r RoutingOperationRow
		if err := rows.Scan(&r.ID, &r.UID, &r.ProductionOrderID, &r.JobID, &r.SequenceNumber, &r.Name, &r.Description, &r.EstimatedDurationHours, &r.MachineID, &r.Notes, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func GetRoutingOperationByUID(pool *pgxpool.Pool, uid string) (RoutingOperationRow, error) {
	ctx := context.Background()
	var r RoutingOperationRow
	err := pool.QueryRow(ctx, `SELECT id, uid, production_order_id, job_id, sequence_number, name, description, estimated_duration_hours, machine_id, notes, created_at, updated_at FROM routing_operations WHERE uid = $1`, uid).Scan(&r.ID, &r.UID, &r.ProductionOrderID, &r.JobID, &r.SequenceNumber, &r.Name, &r.Description, &r.EstimatedDurationHours, &r.MachineID, &r.Notes, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}

func CreateRoutingOperation(pool *pgxpool.Pool, req models.RoutingOperationCreateRequest, productionOrderID *int, jobID *int, machineID *int) (string, error) {
	ctx := context.Background()
	var uid string
	err := pool.QueryRow(ctx, `INSERT INTO routing_operations (production_order_id, job_id, sequence_number, name, description, estimated_duration_hours, machine_id, notes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING uid`, productionOrderID, jobID, req.SequenceNumber, strings.TrimSpace(req.Name), strings.TrimSpace(req.Description), req.EstimatedDurationHours, machineID, strings.TrimSpace(req.Notes)).Scan(&uid)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func UpdateRoutingOperation(pool *pgxpool.Pool, uid string, req models.RoutingOperationUpdateRequest, productionOrderID *int, jobID *int, machineID *int) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1
	if req.ProductionOrderUID != nil {
		setClauses = append(setClauses, fmt.Sprintf("production_order_id = $%d", idx))
		args = append(args, productionOrderID)
		idx++
	}
	if req.JobUID != nil {
		setClauses = append(setClauses, fmt.Sprintf("job_id = $%d", idx))
		args = append(args, jobID)
		idx++
	}
	if req.SequenceNumber != nil {
		setClauses = append(setClauses, fmt.Sprintf("sequence_number = $%d", idx))
		args = append(args, *req.SequenceNumber)
		idx++
	}
	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Name))
		idx++
	}
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Description))
		idx++
	}
	if req.EstimatedDurationHours != nil {
		setClauses = append(setClauses, fmt.Sprintf("estimated_duration_hours = $%d", idx))
		args = append(args, *req.EstimatedDurationHours)
		idx++
	}
	if req.MachineUID != nil {
		setClauses = append(setClauses, fmt.Sprintf("machine_id = $%d", idx))
		args = append(args, machineID)
		idx++
	}
	if req.Notes != nil {
		setClauses = append(setClauses, fmt.Sprintf("notes = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Notes))
		idx++
	}
	if len(setClauses) == 0 {
		return nil
	}
	args = append(args, uid)
	query := fmt.Sprintf("UPDATE routing_operations SET %s, updated_at = NOW() WHERE uid = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeleteRoutingOperation(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM routing_operations WHERE uid = $1`, uid)
	return err
}
