package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SparePartVendorRow struct {
	UID         string
	SupplierUID string
	ModelNo     string
}

func ListSparePartVendors(pool *pgxpool.Pool, spareUID string) ([]SparePartVendorRow, error) {
	query := `
        SELECT spv.uid, v.uid, spv.model_no
        FROM spare_part_vendors spv
        JOIN vendors v ON spv.vendor_id = v.id
        JOIN spare_parts sp ON spv.spare_part_id = sp.id
        WHERE sp.uid = $1
        ORDER BY spv.id ASC
    `
	rows, err := pool.Query(context.Background(), query, spareUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []SparePartVendorRow
	for rows.Next() {
		var r SparePartVendorRow
		if err := rows.Scan(&r.UID, &r.SupplierUID, &r.ModelNo); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func CreateSparePartVendor(pool *pgxpool.Pool, spareID, vendorID int, modelNo string) (string, error) {
	query := `
        INSERT INTO spare_part_vendors (spare_part_id, vendor_id, model_no)
        VALUES ($1, $2, $3)
        RETURNING uid
    `
	var uid string
	err := pool.QueryRow(context.Background(), query, spareID, vendorID, modelNo).Scan(&uid)
	return uid, err
}

func UpdateSparePartVendor(pool *pgxpool.Pool, sparePartID int, vendorID int, modelNo string) (int64, error) {
	cmd, err := pool.Exec(context.Background(),
		`UPDATE spare_part_vendors
         SET model_no = $1
         WHERE spare_part_id = $2 AND vendor_id = $3`,
		modelNo, sparePartID, vendorID,
	)
	if err != nil {
		return 0, err
	}

	return cmd.RowsAffected(), nil
}

func DeleteSparePartVendor(pool *pgxpool.Pool, sparePartID int, vendorID int) (int64, error) {
	cmd, err := pool.Exec(context.Background(),
		`DELETE FROM spare_part_vendors 
         WHERE spare_part_id=$1 AND vendor_id=$2`,
		sparePartID, vendorID,
	)
	if err != nil {
		return 0, err
	}
	return cmd.RowsAffected(), nil
}
