package repo

import (
	"context"
	"fmt"
	"strings"

	"app/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// VendorAddressRow is DB representation for vendor_addresses
type VendorAddressRow struct {
	ID           int
	UID          pgtype.Text
	CreatedAt    pgtype.Timestamp
	CreatedBy    pgtype.Int8
	IsPrimary    bool
	City         pgtype.Text
	State        pgtype.Text
	Country      pgtype.Text
	PinCode      pgtype.Text
	AddressLine1 pgtype.Text
	AddressLine2 pgtype.Text
	AddressLine3 pgtype.Text
	Landmark     pgtype.Text
}

func ListVendorAddresses(pool *pgxpool.Pool, vendorID int) ([]VendorAddressRow, error) {
	ctx := context.Background()
	query := `
		select a.*, va.is_primary
		from address a
		join vendor_addresses va on a.id=va.address_id
		where va.vendor_id = $1 `
	rows, err := pool.Query(ctx, query, vendorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []VendorAddressRow
	for rows.Next() {
		var r VendorAddressRow
		err := rows.Scan(
			&r.ID,
			&r.UID,
			&r.AddressLine1,
			&r.AddressLine2,
			&r.AddressLine3,
			&r.Landmark,
			&r.City,
			&r.State,
			&r.Country,
			&r.PinCode,
			&r.CreatedAt,
			&r.CreatedBy,
			&r.IsPrimary)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func GetVendorAddressByUID(pool *pgxpool.Pool, uid string) (VendorAddressRow, error) {
	ctx := context.Background()
	var r VendorAddressRow
	query := `
		select a.*, va.is_primary
		from address a
		join vendor_addresses va on a.id=va.address_id
		where a.uid = $1
	`
	err := pool.QueryRow(ctx, query, uid).Scan(
		&r.ID,
		&r.UID,
		&r.AddressLine1,
		&r.AddressLine2,
		&r.AddressLine3,
		&r.Landmark,
		&r.City,
		&r.State,
		&r.Country,
		&r.PinCode,
		&r.CreatedAt,
		&r.CreatedBy,
		&r.IsPrimary,
	)
	if err != nil {
		return r, err
	}
	return r, nil
}

func CreateVendorAddress(pool *pgxpool.Pool, req models.VendorAddressCreateRequest, vid int) (string, error) {
	ctx := context.Background()
	var newID int
	err := pool.QueryRow(ctx, `
		INSERT INTO address (address_line_1, address_line_2, address_line_3, landmark, city, state, country, pin_code)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		req.AddressLine1, req.AddressLine2, req.AddressLine3, req.Landmark, req.City, req.State, req.Country, req.PinCode).
		Scan(&newID)
	if err != nil {
		return "0", err
	}

	query := `
		INSERT INTO vendor_addresses(vendor_id, address_id)
		values($1, $2) RETURNING uid
	`
	var newID2 string

	err = pool.QueryRow(ctx, query, vid, newID).Scan(&newID2)

	if err != nil {
		return "0", err
	}

	return newID2, nil
}

// func UpdateVendorAddress(db *sql.DB, id int, req models.VendorAddressUpdateRequest) error {
// 	if req.IsPrimary == nil {
// 		return nil
// 	}
// 	_, err := db.Exec(`UPDATE vendor_addresses SET is_primary = $1 WHERE id = $2`, *req.IsPrimary, id)
// 	return err
// }

func UpdateVendorAddress(pool *pgxpool.Pool, addressID int, req models.VendorAddressUpdateRequest) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []any{}
	idx := 1

	if req.AddressLine1 != nil {
		setClauses = append(setClauses, fmt.Sprintf("address_line_1 = $%d", idx))
		args = append(args, *req.AddressLine1)
		idx++
	}

	if req.AddressLine2 != nil {
		setClauses = append(setClauses, fmt.Sprintf("address_line_2 = $%d", idx))
		args = append(args, *req.AddressLine2)
		idx++
	}

	if req.AddressLine3 != nil {
		setClauses = append(setClauses, fmt.Sprintf("address_line_3 = $%d", idx))
		args = append(args, *req.AddressLine3)
		idx++
	}

	if req.Landmark != nil {
		setClauses = append(setClauses, fmt.Sprintf("landmark = $%d", idx))
		args = append(args, *req.Landmark)
		idx++
	}

	if req.City != nil {
		setClauses = append(setClauses, fmt.Sprintf("city = $%d", idx))
		args = append(args, *req.City)
		idx++
	}

	if req.State != nil {
		setClauses = append(setClauses, fmt.Sprintf("state = $%d", idx))
		args = append(args, *req.State)
		idx++
	}

	if req.Country != nil {
		setClauses = append(setClauses, fmt.Sprintf("country = $%d", idx))
		args = append(args, *req.Country)
		idx++
	}

	if req.PinCode != nil {
		setClauses = append(setClauses, fmt.Sprintf("pin_code = $%d", idx))
		args = append(args, *req.PinCode)
		idx++
	}

	// Nothing to update in address table
	if len(setClauses) > 0 {
		args = append(args, addressID)
		query := "UPDATE address SET " + strings.Join(setClauses, ", ") +
			" WHERE id = " + fmt.Sprintf("$%d", idx)

		_, err := pool.Exec(ctx, query, args...)
		if err != nil {
			return err
		}
	}

	// Update is_primary in vendor_addresses if provided
	if req.IsPrimary != nil {
		_, err := pool.Exec(ctx, `UPDATE vendor_addresses SET is_primary = $1 WHERE address_id = $2`, *req.IsPrimary, addressID)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteVendorAddress(pool *pgxpool.Pool, id int) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM vendor_addresses WHERE id = $1`, id)
	return err
}
