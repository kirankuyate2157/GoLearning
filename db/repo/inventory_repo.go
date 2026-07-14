package repo

import (
	"app/models"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryItemRow struct {
	ID            int
	UID           pgtype.UUID
	PartNumber    pgtype.Text
	Name          string
	Description   pgtype.Text
	Quantity      float64
	Unit          pgtype.Text
	Location      pgtype.Text
	LotNumber     pgtype.Text
	SerialNumber  pgtype.Text
	CreatedAt     pgtype.Timestamp
	UpdatedAt     pgtype.Timestamp
}

type InventoryMovementRow struct {
	ID               int
	UID              pgtype.UUID
	InventoryItemID  int
	InventoryItemName string
	MovementType     string
	Quantity         float64
	ReferenceType    pgtype.Text
	ReferenceUID     pgtype.Text
	Notes            pgtype.Text
	CreatedAt        pgtype.Timestamp
}

type WIPItemRow struct {
	ID                int
	UID               pgtype.UUID
	JobID             int
	InventoryItemID   int
	Quantity          float64
	Status            pgtype.Text
	Notes             pgtype.Text
	CreatedAt         pgtype.Timestamp
	UpdatedAt         pgtype.Timestamp
}

func ListInventoryItems(pool *pgxpool.Pool) ([]InventoryItemRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `SELECT id, uid, part_number, name, description, quantity, unit, location, lot_number, serial_number, created_at, updated_at FROM inventory_items ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []InventoryItemRow
	for rows.Next() {
		var r InventoryItemRow
		if err := rows.Scan(&r.ID, &r.UID, &r.PartNumber, &r.Name, &r.Description, &r.Quantity, &r.Unit, &r.Location, &r.LotNumber, &r.SerialNumber, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func GetInventoryItemByUID(pool *pgxpool.Pool, uid string) (InventoryItemRow, error) {
	ctx := context.Background()
	var r InventoryItemRow
	err := pool.QueryRow(ctx, `SELECT id, uid, part_number, name, description, quantity, unit, location, lot_number, serial_number, created_at, updated_at FROM inventory_items WHERE uid = $1`, uid).Scan(&r.ID, &r.UID, &r.PartNumber, &r.Name, &r.Description, &r.Quantity, &r.Unit, &r.Location, &r.LotNumber, &r.SerialNumber, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}

func CreateInventoryItem(pool *pgxpool.Pool, req models.InventoryItemCreateRequest) (string, error) {
	ctx := context.Background()
	var uid string
	err := pool.QueryRow(ctx, `INSERT INTO inventory_items (part_number, name, description, quantity, unit, location, lot_number, serial_number) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING uid`, strings.TrimSpace(req.PartNumber), strings.TrimSpace(req.Name), strings.TrimSpace(req.Description), req.Quantity, strings.TrimSpace(req.Unit), strings.TrimSpace(req.Location), strings.TrimSpace(req.LotNumber), strings.TrimSpace(req.SerialNumber)).Scan(&uid)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func UpdateInventoryItem(pool *pgxpool.Pool, uid string, req models.InventoryItemUpdateRequest) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1
	if req.PartNumber != nil {
		setClauses = append(setClauses, fmt.Sprintf("part_number = $%d", idx))
		args = append(args, strings.TrimSpace(*req.PartNumber))
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
	if req.Location != nil {
		setClauses = append(setClauses, fmt.Sprintf("location = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Location))
		idx++
	}
	if req.LotNumber != nil {
		setClauses = append(setClauses, fmt.Sprintf("lot_number = $%d", idx))
		args = append(args, strings.TrimSpace(*req.LotNumber))
		idx++
	}
	if req.SerialNumber != nil {
		setClauses = append(setClauses, fmt.Sprintf("serial_number = $%d", idx))
		args = append(args, strings.TrimSpace(*req.SerialNumber))
		idx++
	}
	if len(setClauses) == 0 {
		return nil
	}
	args = append(args, uid)
	query := fmt.Sprintf("UPDATE inventory_items SET %s, updated_at = NOW() WHERE uid = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeleteInventoryItem(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM inventory_items WHERE uid = $1`, uid)
	return err
}

func ListInventoryMovements(pool *pgxpool.Pool) ([]InventoryMovementRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `SELECT im.id, im.uid, im.inventory_item_id, ii.name, im.movement_type, im.quantity, im.reference_type, im.reference_uid, im.notes, im.created_at FROM inventory_movements im JOIN inventory_items ii ON im.inventory_item_id = ii.id ORDER BY im.id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []InventoryMovementRow
	for rows.Next() {
		var r InventoryMovementRow
		if err := rows.Scan(&r.ID, &r.UID, &r.InventoryItemID, &r.InventoryItemName, &r.MovementType, &r.Quantity, &r.ReferenceType, &r.ReferenceUID, &r.Notes, &r.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func CreateInventoryMovement(pool *pgxpool.Pool, req models.InventoryMovementCreateRequest) (string, error) {
	ctx := context.Background()
	var uid string
	err := pool.QueryRow(ctx, `INSERT INTO inventory_movements (inventory_item_id, movement_type, quantity, reference_type, reference_uid, notes) SELECT id, $2, $3, $4, $5, $6 FROM inventory_items WHERE uid = $1 RETURNING uid`, req.InventoryItemUID, strings.TrimSpace(req.MovementType), req.Quantity, strings.TrimSpace(req.ReferenceType), strings.TrimSpace(req.ReferenceUID), strings.TrimSpace(req.Notes)).Scan(&uid)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func ListWIPItems(pool *pgxpool.Pool) ([]WIPItemRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `SELECT id, uid, job_id, inventory_item_id, quantity, status, notes, created_at, updated_at FROM wip_items ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []WIPItemRow
	for rows.Next() {
		var r WIPItemRow
		if err := rows.Scan(&r.ID, &r.UID, &r.JobID, &r.InventoryItemID, &r.Quantity, &r.Status, &r.Notes, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func GetWIPItemByUID(pool *pgxpool.Pool, uid string) (WIPItemRow, error) {
	ctx := context.Background()
	var r WIPItemRow
	err := pool.QueryRow(ctx, `SELECT id, uid, job_id, inventory_item_id, quantity, status, notes, created_at, updated_at FROM wip_items WHERE uid = $1`, uid).Scan(&r.ID, &r.UID, &r.JobID, &r.InventoryItemID, &r.Quantity, &r.Status, &r.Notes, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}

func CreateWIPItem(pool *pgxpool.Pool, req models.WIPItemCreateRequest, jobID int, inventoryItemID int) (string, error) {
	ctx := context.Background()
	var uid string
	err := pool.QueryRow(ctx, `INSERT INTO wip_items (job_id, inventory_item_id, quantity, status, notes) VALUES ($1, $2, $3, $4, $5) RETURNING uid`, jobID, inventoryItemID, req.Quantity, strings.TrimSpace(req.Status), strings.TrimSpace(req.Notes)).Scan(&uid)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func UpdateWIPItem(pool *pgxpool.Pool, uid string, req models.WIPItemUpdateRequest, jobID *int, inventoryItemID *int) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1
	if req.JobUID != nil {
		setClauses = append(setClauses, fmt.Sprintf("job_id = $%d", idx))
		args = append(args, jobID)
		idx++
	}
	if req.InventoryItemUID != nil {
		setClauses = append(setClauses, fmt.Sprintf("inventory_item_id = $%d", idx))
		args = append(args, inventoryItemID)
		idx++
	}
	if req.Quantity != nil {
		setClauses = append(setClauses, fmt.Sprintf("quantity = $%d", idx))
		args = append(args, *req.Quantity)
		idx++
	}
	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Status))
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
	query := fmt.Sprintf("UPDATE wip_items SET %s, updated_at = NOW() WHERE uid = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeleteWIPItem(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM wip_items WHERE uid = $1`, uid)
	return err
}
