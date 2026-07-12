package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"app/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// VendorRow is a lightweight representation of a DB row; controllers will map this to models.Vendor
type VendorRow struct {
	ID        int
	UID       pgtype.Text
	Name      string
	Website   pgtype.Text
	Email     pgtype.Text
	Phone     pgtype.Text
	JoinDate  pgtype.Timestamp
	GSTNo     pgtype.Text
	GSTTypeID pgtype.Int8 // vendors.gst_type (INT)
	GSTType   pgtype.Text // tax.uid (STRING)
	CreatedAt pgtype.Timestamp
}

// ListVendors returns list of vendors from DB
// ListVendors returns vendors, optionally filtered by vendor type NAME
func ListVendors(pool *pgxpool.Pool, typeName string) ([]VendorRow, error) {
	ctx := context.Background()

	var rows pgx.Rows
	var err error

	baseQuery := `
		SELECT
			v.id,
			v.uid,
			v.name,
			v.website,
			v.email,
			v.phone_no,
			v.join_date,
			v.gst_no,
			v.gst_type,
			t.uid AS gst_type,
			v.created_at
		FROM vendors v
		LEFT JOIN tax t ON v.gst_type = t.id
	`

	if strings.TrimSpace(typeName) != "" {
		query := baseQuery + `
			INNER JOIN vendor_types vt ON vt.vendor_id = v.id
			INNER JOIN types ttype ON vt.type_id = ttype.id
			WHERE LOWER(ttype.name) = LOWER($1)
			ORDER BY v.id ASC;
		`
		rows, err = pool.Query(ctx, query, typeName)
	} else {
		rows, err = pool.Query(ctx, baseQuery+` ORDER BY v.id ASC;`)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []VendorRow
	for rows.Next() {
		var r VendorRow
		err := rows.Scan(
			&r.ID,
			&r.UID,
			&r.Name,
			&r.Website,
			&r.Email,
			&r.Phone,
			&r.JoinDate,
			&r.GSTNo,
			&r.GSTTypeID,
			&r.GSTType,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}

	return out, rows.Err()
}

// GetVendorByUID returns vendor by ID
func GetVendorByUID(pool *pgxpool.Pool, uid string) (VendorRow, error) {
	ctx := context.Background()
	var r VendorRow

	err := pool.QueryRow(ctx, `
	SELECT
		v.id,
		v.uid,
		v.name,
		v.website,
		v.email,
		v.phone_no,
		v.join_date,
		v.gst_no,
		v.gst_type,
		t.uid AS gst_type,
		v.created_at
	FROM vendors v
	LEFT JOIN tax t ON v.gst_type = t.id
	WHERE v.uid = $1
`, uid).Scan(
		&r.ID,
		&r.UID,
		&r.Name,
		&r.Website,
		&r.Email,
		&r.Phone,
		&r.JoinDate,
		&r.GSTNo,
		&r.GSTTypeID,
		&r.GSTType,
		&r.CreatedAt,
	)

	return r, err
}

// CreateVendor inserts a vendor and returns new id
func CreateVendor(pool *pgxpool.Pool, req models.VendorCreateRequest, taxID int) (string, error) {
	ctx := context.Background()
	var joinDate interface{}
	if req.JoinDate != "" {
		if t, err := time.Parse("2006-01-02", req.JoinDate); err == nil {
			joinDate = t
		}
	}
	var newID string
	// address_id and contact_person_id are optional and not included in create request currently
	// err := pool.QueryRow(ctx, `INSERT INTO vendors (uid, name, website, email, phone_no, join_date, gst_no, tax_id, created_at) VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, NOW()) RETURNING uid`, req.Name, req.Website, req.Email, req.PhoneNo, joinDate, req.GSTNo, taxID).Scan(&newID)

	err := pool.QueryRow(ctx, `
	INSERT INTO vendors
	(uid, name, website, email, phone_no, join_date, gst_no, gst_type, created_at)
	VALUES (gen_random_uuid(), $1,$2,$3,$4,$5,$6,$7,NOW())
	RETURNING uid
`,
		req.Name,
		req.Website,
		req.Email,
		req.PhoneNo,
		joinDate,
		req.GSTNo,
		taxID,
	).Scan(&newID)

	if err != nil {
		return "0", err
	}
	return newID, nil
}

// UpdateVendor updates a vendor based on provided fields
func UpdateVendor(pool *pgxpool.Pool, uid string, req models.VendorUpdateRequest, taxID *int) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1
	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, *req.Name)
		idx++
	}
	if req.Website != nil {
		setClauses = append(setClauses, fmt.Sprintf("website = $%d", idx))
		args = append(args, *req.Website)
		idx++
	}
	if req.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", idx))
		args = append(args, *req.Email)
		idx++
	}
	if req.PhoneNo != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone_no = $%d", idx))
		args = append(args, *req.PhoneNo)
		idx++
	}
	if req.JoinDate != nil {
		if t, err := time.Parse("2006-01-02", *req.JoinDate); err == nil {
			setClauses = append(setClauses, fmt.Sprintf("join_date = $%d", idx))
			args = append(args, t)
			idx++
		}
	}
	if req.GSTNo != nil {
		setClauses = append(setClauses, fmt.Sprintf("gst_no = $%d", idx))
		args = append(args, *req.GSTNo)
		idx++
	}
	if taxID != nil {
		setClauses = append(setClauses, fmt.Sprintf("gst_type = $%d", idx))
		args = append(args, *taxID)
		idx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, uid)
	query := fmt.Sprintf("UPDATE vendors SET %s WHERE uid = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

// DeleteVendor removes a vendor by UID
func DeleteVendor(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM vendors WHERE uid = $1`, uid)
	return err
}
