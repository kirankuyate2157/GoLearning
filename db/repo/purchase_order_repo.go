package repo

import (
	"app/models"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PurchaseOrderRow struct {
	ID                   int
	UID                  pgtype.UUID
	VendorUID            pgtype.UUID
	VendorName           pgtype.Text
	OrderNumber          string
	Status               pgtype.Text
	OrderDate            pgtype.Date
	ExpectedDeliveryDate pgtype.Date
	TotalAmount          float64
	Currency             pgtype.Text
	Notes                pgtype.Text
	CreatedAt            pgtype.Timestamp
	UpdatedAt            pgtype.Timestamp
}

type PurchaseOrderLineRow struct {
	ID               int
	UID              pgtype.UUID
	PurchaseOrderUID pgtype.UUID
	InventoryItemUID pgtype.UUID
	ItemDescription  pgtype.Text
	Quantity         float64
	UnitPrice        float64
	TotalPrice       float64
	Notes            pgtype.Text
	CreatedAt        pgtype.Timestamp
	UpdatedAt        pgtype.Timestamp
}

func ListPurchaseOrders(pool *pgxpool.Pool) ([]PurchaseOrderRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `SELECT po.id, po.uid, v.uid AS vendor_uid, v.name AS vendor_name, po.order_number, po.status, po.order_date, po.expected_delivery_date, po.total_amount::float8 AS total_amount, po.currency, po.notes, po.created_at, po.updated_at FROM purchase_orders po LEFT JOIN vendors v ON po.vendor_id = v.id ORDER BY po.id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []PurchaseOrderRow
	for rows.Next() {
		var r PurchaseOrderRow
		if err := rows.Scan(&r.ID, &r.UID, &r.VendorUID, &r.VendorName, &r.OrderNumber, &r.Status, &r.OrderDate, &r.ExpectedDeliveryDate, &r.TotalAmount, &r.Currency, &r.Notes, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func ListPurchaseOrdersByVendorUID(pool *pgxpool.Pool, vendorUID string) ([]PurchaseOrderRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `SELECT po.id, po.uid, v.uid AS vendor_uid, v.name AS vendor_name, po.order_number, po.status, po.order_date, po.expected_delivery_date, po.total_amount::float8 AS total_amount, po.currency, po.notes, po.created_at, po.updated_at FROM purchase_orders po JOIN vendors v ON po.vendor_id = v.id WHERE v.uid = $1 ORDER BY po.id ASC`, vendorUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []PurchaseOrderRow
	for rows.Next() {
		var r PurchaseOrderRow
		if err := rows.Scan(&r.ID, &r.UID, &r.VendorUID, &r.VendorName, &r.OrderNumber, &r.Status, &r.OrderDate, &r.ExpectedDeliveryDate, &r.TotalAmount, &r.Currency, &r.Notes, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func GetPurchaseOrderByUID(pool *pgxpool.Pool, uid string) (PurchaseOrderRow, error) {
	ctx := context.Background()
	var r PurchaseOrderRow
	err := pool.QueryRow(ctx, `SELECT po.id, po.uid, v.uid AS vendor_uid, v.name AS vendor_name, po.order_number, po.status, po.order_date, po.expected_delivery_date, po.total_amount::float8 AS total_amount, po.currency, po.notes, po.created_at, po.updated_at FROM purchase_orders po LEFT JOIN vendors v ON po.vendor_id = v.id WHERE po.uid = $1`, uid).Scan(&r.ID, &r.UID, &r.VendorUID, &r.VendorName, &r.OrderNumber, &r.Status, &r.OrderDate, &r.ExpectedDeliveryDate, &r.TotalAmount, &r.Currency, &r.Notes, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}

func CreatePurchaseOrder(pool *pgxpool.Pool, req models.PurchaseOrderCreateRequest, vendorID int) (string, error) {
	ctx := context.Background()
	orderDate, err := nullableDate(req.OrderDate)
	if err != nil {
		return "", err
	}
	expectedDate, err := nullableDate(req.ExpectedDeliveryDate)
	if err != nil {
		return "", err
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "draft"
	}

	var uid string
	err = pool.QueryRow(ctx, `INSERT INTO purchase_orders (vendor_id, order_number, status, order_date, expected_delivery_date, total_amount, currency, notes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING uid`, vendorID, strings.TrimSpace(req.OrderNumber), status, orderDate, expectedDate, req.TotalAmount, strings.TrimSpace(req.Currency), strings.TrimSpace(req.Notes)).Scan(&uid)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func UpdatePurchaseOrder(pool *pgxpool.Pool, uid string, req models.PurchaseOrderUpdateRequest, vendorID *int) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	if req.VendorUID != nil {
		setClauses = append(setClauses, fmt.Sprintf("vendor_id = $%d", idx))
		args = append(args, vendorID)
		idx++
	}
	if req.OrderNumber != nil {
		setClauses = append(setClauses, fmt.Sprintf("order_number = $%d", idx))
		args = append(args, strings.TrimSpace(*req.OrderNumber))
		idx++
	}
	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Status))
		idx++
	}
	if req.OrderDate != nil {
		dateArg, err := nullableDate(*req.OrderDate)
		if err != nil {
			return err
		}
		setClauses = append(setClauses, fmt.Sprintf("order_date = $%d", idx))
		args = append(args, dateArg)
		idx++
	}
	if req.ExpectedDeliveryDate != nil {
		dateArg, err := nullableDate(*req.ExpectedDeliveryDate)
		if err != nil {
			return err
		}
		setClauses = append(setClauses, fmt.Sprintf("expected_delivery_date = $%d", idx))
		args = append(args, dateArg)
		idx++
	}
	if req.TotalAmount != nil {
		setClauses = append(setClauses, fmt.Sprintf("total_amount = $%d", idx))
		args = append(args, *req.TotalAmount)
		idx++
	}
	if req.Currency != nil {
		setClauses = append(setClauses, fmt.Sprintf("currency = $%d", idx))
		args = append(args, strings.TrimSpace(*req.Currency))
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
	query := fmt.Sprintf("UPDATE purchase_orders SET %s, updated_at = NOW() WHERE uid = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeletePurchaseOrder(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM purchase_orders WHERE uid = $1`, uid)
	return err
}

func ListPurchaseOrderLines(pool *pgxpool.Pool, purchaseOrderUID string) ([]PurchaseOrderLineRow, error) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, `SELECT pol.id, pol.uid, po.uid AS purchase_order_uid, ii.uid AS inventory_item_uid, pol.item_description, pol.quantity::float8 AS quantity, pol.unit_price::float8 AS unit_price, pol.total_price::float8 AS total_price, pol.notes, pol.created_at, pol.updated_at FROM purchase_order_lines pol JOIN purchase_orders po ON pol.purchase_order_id = po.id LEFT JOIN inventory_items ii ON pol.inventory_item_id = ii.id WHERE po.uid = $1 ORDER BY pol.id ASC`, purchaseOrderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []PurchaseOrderLineRow
	for rows.Next() {
		var r PurchaseOrderLineRow
		if err := rows.Scan(&r.ID, &r.UID, &r.PurchaseOrderUID, &r.InventoryItemUID, &r.ItemDescription, &r.Quantity, &r.UnitPrice, &r.TotalPrice, &r.Notes, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func CreatePurchaseOrderLine(pool *pgxpool.Pool, req models.PurchaseOrderLineCreateRequest, purchaseOrderID int, inventoryItemID *int) (string, error) {
	ctx := context.Background()
	var uid string
	query := `INSERT INTO purchase_order_lines (purchase_order_id, inventory_item_id, item_description, quantity, unit_price, total_price, notes) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING uid`
	if inventoryItemID != nil {
		if err := pool.QueryRow(ctx, query, purchaseOrderID, inventoryItemID, strings.TrimSpace(req.ItemDescription), req.Quantity, req.UnitPrice, req.TotalPrice, strings.TrimSpace(req.Notes)).Scan(&uid); err != nil {
			return "", err
		}
	} else {
		if err := pool.QueryRow(ctx, query, purchaseOrderID, nil, strings.TrimSpace(req.ItemDescription), req.Quantity, req.UnitPrice, req.TotalPrice, strings.TrimSpace(req.Notes)).Scan(&uid); err != nil {
			return "", err
		}
	}
	return uid, nil
}

func UpdatePurchaseOrderLine(pool *pgxpool.Pool, uid string, req models.PurchaseOrderLineUpdateRequest, inventoryItemID *int) error {
	ctx := context.Background()
	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	if req.InventoryItemUID != nil {
		setClauses = append(setClauses, fmt.Sprintf("inventory_item_id = $%d", idx))
		args = append(args, inventoryItemID)
		idx++
	}
	if req.ItemDescription != nil {
		setClauses = append(setClauses, fmt.Sprintf("item_description = $%d", idx))
		args = append(args, strings.TrimSpace(*req.ItemDescription))
		idx++
	}
	if req.Quantity != nil {
		setClauses = append(setClauses, fmt.Sprintf("quantity = $%d", idx))
		args = append(args, *req.Quantity)
		idx++
	}
	if req.UnitPrice != nil {
		setClauses = append(setClauses, fmt.Sprintf("unit_price = $%d", idx))
		args = append(args, *req.UnitPrice)
		idx++
	}
	if req.TotalPrice != nil {
		setClauses = append(setClauses, fmt.Sprintf("total_price = $%d", idx))
		args = append(args, *req.TotalPrice)
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
	query := fmt.Sprintf("UPDATE purchase_order_lines SET %s, updated_at = NOW() WHERE uid = $%d", strings.Join(setClauses, ", "), idx)
	_, err := pool.Exec(ctx, query, args...)
	return err
}

func DeletePurchaseOrderLine(pool *pgxpool.Pool, uid string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `DELETE FROM purchase_order_lines WHERE uid = $1`, uid)
	return err
}
