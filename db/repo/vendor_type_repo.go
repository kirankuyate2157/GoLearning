package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Type model for returning vendor type details
type Type struct {
	ID          int
	UID         string `json:"uid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type VendorTypeRow struct {
	TypeUID         pgtype.Text
	VendorAssignUID pgtype.Text
	Name            pgtype.Text
	Description     pgtype.Text
	Category        pgtype.Text
}

// CreateVendorType assigns a vendor type to a vendor
func CreateVendorType(pool *pgxpool.Pool, vendorUID, vendorTypeUID string) (map[string]string, error) {
	ctx := context.Background()
	var vendorID, typeID int

	// Resolve vendor_id
	if err := pool.QueryRow(ctx, `SELECT id FROM vendors WHERE uid = $1`, vendorUID).Scan(&vendorID); err != nil {
		return nil, fmt.Errorf("vendor not found: %v", err)
	}

	// Resolve type_id
	if err := pool.QueryRow(ctx, `SELECT id FROM types WHERE uid = $1`, vendorTypeUID).Scan(&typeID); err != nil {
		return nil, fmt.Errorf("type not found: %v", err)
	}

	// Insert record (created_by defaults 0, created_at auto)
	var uid string
	query := `
		INSERT INTO vendor_types (vendor_id, type_id)
		VALUES ($1, $2)
		RETURNING uid;
	`
	if err := pool.QueryRow(ctx, query, vendorID, typeID).Scan(&uid); err != nil {
		return nil, fmt.Errorf("insert failed: %v", err)
	}

	return map[string]string{
		"vendor_uid":      vendorUID,
		"vendor_type_uid": vendorTypeUID,
		"uid":             uid,
	}, nil
}

// GetVendorTypes retrieves all type details assigned to a vendor
func GetVendorTypes(pool *pgxpool.Pool, vendorUID string) ([]VendorTypeRow, error) {
	ctx := context.Background()
	query := `
		SELECT 
			t.uid AS type_uid,
			vt.uid AS vendor_assign_uid,
			t.name,
			t.description,
			t.category
		FROM vendor_types vt
		INNER JOIN vendors v ON vt.vendor_id = v.id
		INNER JOIN types t ON vt.type_id = t.id
		WHERE v.uid = $1
		ORDER BY t.name;
	`

	rows, err := pool.Query(ctx, query, vendorUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []VendorTypeRow

	for rows.Next() {
		var r VendorTypeRow

		if err := rows.Scan(
			&r.TypeUID,
			&r.VendorAssignUID,
			&r.Name,
			&r.Description,
			&r.Category,
		); err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteVendorType removes a vendor type assignment by vendorAssign_uid
func DeleteVendorType(pool *pgxpool.Pool, vendorUID, vendorAssignUID string) error {
	ctx := context.Background()
	var vendorID int

	// Ensure vendor exists and resolve id
	if err := pool.QueryRow(ctx, `SELECT id FROM vendors WHERE uid = $1`, vendorUID).Scan(&vendorID); err != nil {
		return fmt.Errorf("vendor not found: %v", err)
	}

	// Delete by vendor_types.uid (vendorAssignUID)
	query := `DELETE FROM vendor_types WHERE vendor_id = $1 AND uid = $2`
	result, err := pool.Exec(ctx, query, vendorID, vendorAssignUID)
	if err != nil {
		return err
	}

	count := result.RowsAffected()
	if count == 0 {
		return fmt.Errorf("no record found to delete")
	}

	return nil
}
