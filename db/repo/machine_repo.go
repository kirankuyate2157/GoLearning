package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"app/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MachineRow represents a machine row fetched from the database
type MachineRow struct {
	ID                    int
	UID                   pgtype.Text
	Name                  string
	Code                  string
	Description           pgtype.Text
	Model                 pgtype.Text
	Vendor                pgtype.Int8
	SerialNo              pgtype.Text
	InstallationDate      pgtype.Timestamp
	CreatedAt             pgtype.Timestamp
	CreatedBy             pgtype.Int8
	VendorUID             pgtype.Text
	TypeID                pgtype.Text
	MinOperationValue     pgtype.Float8
	MinOperationValueUnit pgtype.Text
	MaxOperationValue     pgtype.Float8
	MaxOperationValueUnit pgtype.Text
	ElectricityUPH        pgtype.Float8
	RentPurchaseCost      pgtype.Float8
	ConsumableCost        pgtype.Float8
	WagesPerHour          pgtype.Float8
}

// ListMachines returns all machines ordered by ID

func ListMachines(pool *pgxpool.Pool) ([]MachineRow, error) {
	ctx := context.Background()

	query := `
		SELECT 
			m.id,
			m.uid,
			m.name,
			m.code,
			m.description,
			m.model,
			v.uid AS vendor_uid,
			m.serial_no,
			m.installation_date,
			m.created_at,
			m.created_by,
			t.uid AS type_uid,
			m.min_operation_value,
			u1.uid AS min_op_unit_uid,
			m.max_operation_value,
			u2.uid AS max_op_unit_uid,
			m.electrity_uph,
			m.rent_purchase_cost,
			m.consumable_cost,
			m.wages_per_hour
			FROM machines m
			LEFT JOIN vendors v ON m.vendor = v.id
			LEFT JOIN types t  ON m.type_id = t.id
			LEFT JOIN types u1 ON m.min_operation_value_unit = u1.id
			LEFT JOIN types u2 ON m.max_operation_value_unit = u2.id

		ORDER BY m.id ASC;
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []MachineRow

	for rows.Next() {
		var r MachineRow

		err := rows.Scan(
			&r.ID,
			&r.UID,
			&r.Name,
			&r.Code,
			&r.Description,
			&r.Model,

			&r.VendorUID,

			&r.SerialNo,
			&r.InstallationDate,
			&r.CreatedAt,
			&r.CreatedBy,

			&r.TypeID,
			&r.MinOperationValue,
			&r.MinOperationValueUnit,
			&r.MaxOperationValue,
			&r.MaxOperationValueUnit,

			&r.ElectricityUPH,
			&r.RentPurchaseCost,
			&r.ConsumableCost,
			&r.WagesPerHour,
		)

		if err != nil {
			return nil, err
		}

		out = append(out, r)
	}

	return out, nil
}

// GetMachineByUID returns a single machine by its UID
func GetMachineByUID(pool *pgxpool.Pool, uid string) (MachineRow, error) {
	ctx := context.Background()

	query := `
		SELECT 
			m.id,
			m.uid,
			m.name,
			m.code,
			m.description,
			m.model,
			v.uid AS vendor_uid,
			m.serial_no,
			m.installation_date,
			m.created_at,
			m.created_by,
			t.uid  AS type_uid,
			m.min_operation_value,
			u1.uid AS min_op_unit_uid,
			m.max_operation_value,
			u2.uid AS max_op_unit_uid,
			m.electrity_uph,
			m.rent_purchase_cost,
			m.consumable_cost,
			m.wages_per_hour
			FROM machines m
			LEFT JOIN vendors v ON m.vendor = v.id
			LEFT JOIN types   t ON m.type_id = t.id
			LEFT JOIN types  u1 ON m.min_operation_value_unit = u1.id
			LEFT JOIN types  u2 ON m.max_operation_value_unit = u2.id
			WHERE m.uid = $1
	`

	var r MachineRow

	err := pool.QueryRow(ctx, query, uid).Scan(
		&r.ID,
		&r.UID,
		&r.Name,
		&r.Code,
		&r.Description,
		&r.Model,

		&r.VendorUID,

		&r.SerialNo,
		&r.InstallationDate,
		&r.CreatedAt,
		&r.CreatedBy,

		&r.TypeID, // now contains UID
		&r.MinOperationValue,
		&r.MinOperationValueUnit, // now contains UID
		&r.MaxOperationValue,
		&r.MaxOperationValueUnit, // now contains UID

		&r.ElectricityUPH,
		&r.RentPurchaseCost,
		&r.ConsumableCost,
		&r.WagesPerHour,
	)

	return r, err
}

// CreateMachine inserts a new machine and returns its UID
func CreateMachine(pool *pgxpool.Pool, req models.MachineCreateRequest,
	vendorID, typeID, minUnitID, maxUnitID *int) (string, error) {

	ctx := context.Background()
	var newUID string

	err := pool.QueryRow(ctx,
		`INSERT INTO machines (
			name, code, description, model, vendor, serial_no, installation_date,
			type_id, min_operation_value, min_operation_value_unit,
			max_operation_value, max_operation_value_unit,
			electrity_uph, rent_purchase_cost, consumable_cost, wages_per_hour
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7,
				$8, $9, $10,
				$11, $12,
				$13, $14, $15, $16)
		RETURNING uid`,
		req.Name,
		req.Code,
		nullIfEmpty(req.Description),
		nullIfEmpty(req.Model),
		vendorID,
		nullIfEmpty(req.SerialNo),
		nullIfEmpty(req.InstallationDate),
		typeID,
		req.MinOperationValue,
		minUnitID,
		req.MaxOperationValue,
		maxUnitID,
		req.ElectricityUPH,
		req.RentPurchaseCost,
		req.ConsumableCost,
		req.WagesPerHour,
	).Scan(&newUID)

	if err != nil {
		return "", err
	}

	return newUID, nil
}

// UpdateMachine applies partial updates to a machine identified by UID
func UpdateMachine(pool *pgxpool.Pool, uid string, req models.MachineUpdateRequest,
	vendorID, typeID, minUnitID, maxUnitID *int,
) error {

	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, *req.Name)
		idx++
	}

	if req.Code != nil {
		setClauses = append(setClauses, fmt.Sprintf("code = $%d", idx))
		args = append(args, *req.Code)
		idx++
	}

	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", idx))
		args = append(args, nullIfEmpty(*req.Description))
		idx++
	}

	if req.Model != nil {
		setClauses = append(setClauses, fmt.Sprintf("model = $%d", idx))
		args = append(args, nullIfEmpty(*req.Model))
		idx++
	}

	// vendor UID → vendor ID
	if vendorID != nil {
		setClauses = append(setClauses, fmt.Sprintf("vendor = $%d", idx))
		args = append(args, *vendorID)
		idx++
	}

	if req.SerialNo != nil {
		setClauses = append(setClauses, fmt.Sprintf("serial_no = $%d", idx))
		args = append(args, nullIfEmpty(*req.SerialNo))
		idx++
	}

	if req.InstallationDate != nil {
		if strings.TrimSpace(*req.InstallationDate) == "" {
			setClauses = append(setClauses, "installation_date = NULL")
		} else if t, err := time.Parse("2006-01-02", *req.InstallationDate); err == nil {
			setClauses = append(setClauses, fmt.Sprintf("installation_date = $%d", idx))
			args = append(args, t)
			idx++
		}
	}

	// type UID → type ID
	if typeID != nil {
		setClauses = append(setClauses, fmt.Sprintf("type_id = $%d", idx))
		args = append(args, *typeID)
		idx++
	}

	if req.MinOperationValue != nil {
		setClauses = append(setClauses, fmt.Sprintf("min_operation_value = $%d", idx))
		args = append(args, *req.MinOperationValue)
		idx++
	}

	// min unit UID → min unit ID
	if minUnitID != nil {
		setClauses = append(setClauses, fmt.Sprintf("min_operation_value_unit = $%d", idx))
		args = append(args, *minUnitID)
		idx++
	}

	if req.MaxOperationValue != nil {
		setClauses = append(setClauses, fmt.Sprintf("max_operation_value = $%d", idx))
		args = append(args, *req.MaxOperationValue)
		idx++
	}

	// max unit UID → max unit ID
	if maxUnitID != nil {
		setClauses = append(setClauses, fmt.Sprintf("max_operation_value_unit = $%d", idx))
		args = append(args, *maxUnitID)
		idx++
	}

	if req.ElectricityUPH != nil {
		setClauses = append(setClauses, fmt.Sprintf("electrity_uph = $%d", idx))
		args = append(args, *req.ElectricityUPH)
		idx++
	}

	if req.RentPurchaseCost != nil {
		setClauses = append(setClauses, fmt.Sprintf("rent_purchase_cost = $%d", idx))
		args = append(args, *req.RentPurchaseCost)
		idx++
	}

	if req.ConsumableCost != nil {
		setClauses = append(setClauses, fmt.Sprintf("consumable_cost = $%d", idx))
		args = append(args, *req.ConsumableCost)
		idx++
	}

	if req.WagesPerHour != nil {
		setClauses = append(setClauses, fmt.Sprintf("wages_per_hour = $%d", idx))
		args = append(args, *req.WagesPerHour)
		idx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, uid)
	query := fmt.Sprintf("UPDATE machines SET %s WHERE uid = $%d", strings.Join(setClauses, ", "), idx)

	_, err := pool.Exec(ctx, query, args...)
	return err
}

// DeleteMachine removes a machine by UID
func DeleteMachine(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM machines WHERE uid = $1`, uid)
	return err
}

// nullIfEmpty returns nil if the provided string is empty after trimming
