package repo

import (
	"context"
	"fmt"
	"strings"

	"app/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// VendorContactRow is a lightweight representation of a DB row; controllers will map this to models.VendorContact
type VendorContactRow struct {
	ID        int
	UID       pgtype.Text
	Name      pgtype.Text
	Email     pgtype.Text
	Mobile    pgtype.Text
	CreatedAt pgtype.Timestamp
}

// ListVendorContacts returns list of contacts for a vendor from DB
func ListVendorContacts(pool *pgxpool.Pool, vendorID int) ([]VendorContactRow, error) {
	ctx := context.Background()
	query := `
		SELECT p.id, p.uid, p.name, p.email, p.mobile_number, p.created_at 
		FROM vendor_contacts vp 
		JOIN person p ON vp.contact_id = p.id
		WHERE vp.vendor_id = $1 ORDER BY p.id ASC
	`
	rows, err := pool.Query(ctx, query, vendorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []VendorContactRow
	for rows.Next() {
		var r VendorContactRow
		if err := rows.Scan(&r.ID, &r.UID, &r.Name, &r.Email, &r.Mobile, &r.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateVendorContact inserts a vendor contact and returns new id
func CreateVendorContact(pool *pgxpool.Pool, req models.VendorContactCreateRequest, vendorID int) (string, error) {
	ctx := context.Background()
	var newID int
	query := `
		INSERT INTO person (name, email, mobile_number) 
		VALUES ($1, $2, $3) RETURNING id
	`
	err := pool.QueryRow(ctx, query, req.Name, req.Email, req.Mobile).Scan(&newID)
	if err != nil {
		return "0", err
	}

	var newUID2 string

	query = `
		INSERT INTO vendor_contacts(vendor_id, contact_id)
		VALUES ($1, $2)
		RETURNING uid
	`

	err = pool.QueryRow(ctx, query, vendorID, newID).Scan(&newUID2)

	if err != nil {
		return "0", err
	}

	return newUID2, nil
}

// GetVendorContactByUID returns vendor contact by UID
func GetVendorContactByUID(pool *pgxpool.Pool, uid string) (VendorContactRow, error) {
	ctx := context.Background()
	var r VendorContactRow
	query := `
		SELECT p.id, p.uid, p.name, p.email, p.mobile_number, p.created_at 
		FROM vendor_contacts vp 
		JOIN person p ON vp.contact_id = p.id
		WHERE p.uid = $1 ORDER BY p.id ASC
	`
	err := pool.QueryRow(ctx, query, uid).Scan(&r.ID, &r.UID, &r.Name, &r.Email, &r.Mobile, &r.CreatedAt)
	if err != nil {
		return r, err
	}
	return r, nil
}

// UpdateVendorContact updates a vendor contact by person ID
func UpdateVendorContact(pool *pgxpool.Pool, personID int, req models.VendorContactUpdateRequest) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []any{}
	idx := 1

	// ---- NAME ----
	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, *req.Name)
		idx++
	}

	// ---- EMAIL ----
	if req.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", idx))
		args = append(args, *req.Email)
		idx++
	}

	// ---- MOBILE ----
	if req.Mobile != nil {
		setClauses = append(setClauses, fmt.Sprintf("mobile_number = $%d", idx))
		args = append(args, *req.Mobile)
		idx++
	}

	// Nothing to update
	if len(setClauses) == 0 {
		return nil
	}

	// WHERE clause
	args = append(args, personID)
	query := "UPDATE person SET " + strings.Join(setClauses, ", ") +
		" WHERE id = " + fmt.Sprintf("$%d", idx)

	_, err := pool.Exec(ctx, query, args...)
	return err
}

// DeleteVendorContact deletes a vendor contact by id
func DeleteVendorContact(pool *pgxpool.Pool, id int) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM public.vendor_contacts WHERE contact_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = pool.Exec(ctx, `DELETE FROM person WHERE id = $1`, id)
	return err
}
