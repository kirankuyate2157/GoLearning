package controllers

import (
	dbrepo "app/db/repo"
	"app/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func GetPurchaseOrders(c echo.Context) error {
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.PurchaseOrder{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	rows, err := dbrepo.ListPurchaseOrders(pool)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query purchase orders: %v", err)})
	}

	out := make([]models.PurchaseOrder, 0, len(rows))
	for _, r := range rows {
		po := models.PurchaseOrder{UID: uuidString(r.UID), VendorUID: uuidString(r.VendorUID), VendorName: r.VendorName.String, OrderNumber: r.OrderNumber, TotalAmount: r.TotalAmount}
		if r.Status.Valid {
			po.Status = r.Status.String
		}
		if r.OrderDate.Valid {
			po.OrderDate = r.OrderDate.Time.Format("2006-01-02")
		}
		if r.ExpectedDeliveryDate.Valid {
			po.ExpectedDeliveryDate = r.ExpectedDeliveryDate.Time.Format("2006-01-02")
		}
		if r.Currency.Valid {
			po.Currency = r.Currency.String
		}
		if r.Notes.Valid {
			po.Notes = r.Notes.String
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			po.CreatedAt = &t
		}
		if r.UpdatedAt.Valid {
			t := r.UpdatedAt.Time
			po.UpdatedAt = &t
		}
		out = append(out, po)
	}
	return c.JSON(http.StatusOK, out)
}

func GetPurchaseOrder(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	r, err := dbrepo.GetPurchaseOrderByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query purchase order: %v", err)})
	}

	po := models.PurchaseOrder{UID: uuidString(r.UID), VendorUID: uuidString(r.VendorUID), VendorName: r.VendorName.String, OrderNumber: r.OrderNumber, TotalAmount: r.TotalAmount}
	if r.Status.Valid {
		po.Status = r.Status.String
	}
	if r.OrderDate.Valid {
		po.OrderDate = r.OrderDate.Time.Format("2006-01-02")
	}
	if r.ExpectedDeliveryDate.Valid {
		po.ExpectedDeliveryDate = r.ExpectedDeliveryDate.Time.Format("2006-01-02")
	}
	if r.Currency.Valid {
		po.Currency = r.Currency.String
	}
	if r.Notes.Valid {
		po.Notes = r.Notes.String
	}
	if r.CreatedAt.Valid {
		t := r.CreatedAt.Time
		po.CreatedAt = &t
	}
	if r.UpdatedAt.Valid {
		t := r.UpdatedAt.Time
		po.UpdatedAt = &t
	}
	return c.JSON(http.StatusOK, po)
}

func CreatePurchaseOrder(c echo.Context) error {
	var req models.PurchaseOrderCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if strings.TrimSpace(req.VendorUID) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "vendor_uid is required"})
	}
	vendorUID := strings.TrimSpace(req.VendorUID)

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "db not configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	vendor, err := dbrepo.GetVendorByUID(pool, vendorUID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "vendor uid not found"})
	}

	uid, err := dbrepo.CreatePurchaseOrder(pool, req, vendor.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": uid})
}

func UpdatePurchaseOrder(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var req models.PurchaseOrderUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "db not configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	var vendorID *int
	if req.VendorUID != nil && strings.TrimSpace(*req.VendorUID) != "" {
		vendor, err := dbrepo.GetVendorByUID(pool, strings.TrimSpace(*req.VendorUID))
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "vendor uid not found"})
		}
		vendorID = &vendor.ID
	}

	if err := dbrepo.UpdatePurchaseOrder(pool, uid, req, vendorID); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

func DeletePurchaseOrder(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	if err := dbrepo.DeletePurchaseOrder(pool, uid); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete purchase order: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}

func GetVendorPurchaseOrders(c echo.Context) error {
	vendorUID := strings.TrimSpace(c.Param("vendor_uid"))
	if vendorUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid vendor uid"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.PurchaseOrder{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	if _, err := dbrepo.GetVendorByUID(pool, vendorUID); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "vendor not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("vendor lookup error: %v", err)})
	}

	rows, err := dbrepo.ListPurchaseOrdersByVendorUID(pool, vendorUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query vendor purchase orders: %v", err)})
	}

	out := make([]models.PurchaseOrder, 0, len(rows))
	for _, r := range rows {
		po := models.PurchaseOrder{UID: uuidString(r.UID), VendorUID: uuidString(r.VendorUID), VendorName: r.VendorName.String, OrderNumber: r.OrderNumber, TotalAmount: r.TotalAmount}
		if r.Status.Valid {
			po.Status = r.Status.String
		}
		if r.OrderDate.Valid {
			po.OrderDate = r.OrderDate.Time.Format("2006-01-02")
		}
		if r.ExpectedDeliveryDate.Valid {
			po.ExpectedDeliveryDate = r.ExpectedDeliveryDate.Time.Format("2006-01-02")
		}
		if r.Currency.Valid {
			po.Currency = r.Currency.String
		}
		if r.Notes.Valid {
			po.Notes = r.Notes.String
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			po.CreatedAt = &t
		}
		if r.UpdatedAt.Valid {
			t := r.UpdatedAt.Time
			po.UpdatedAt = &t
		}
		out = append(out, po)
	}
	return c.JSON(http.StatusOK, out)
}

func GetPurchaseOrderLines(c echo.Context) error {
	purchaseUID := strings.TrimSpace(c.Param("uid"))
	if purchaseUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid purchase order uid"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.PurchaseOrderLine{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	if _, err := dbrepo.GetPurchaseOrderByUID(pool, purchaseUID); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "purchase order not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("purchase order lookup error: %v", err)})
	}

	rows, err := dbrepo.ListPurchaseOrderLines(pool, purchaseUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query purchase order lines: %v", err)})
	}

	out := make([]models.PurchaseOrderLine, 0, len(rows))
	for _, r := range rows {
		line := models.PurchaseOrderLine{UID: uuidString(r.UID), PurchaseOrderUID: uuidString(r.PurchaseOrderUID), InventoryItemUID: uuidString(r.InventoryItemUID), ItemDescription: r.ItemDescription.String, Quantity: r.Quantity, UnitPrice: r.UnitPrice, TotalPrice: r.TotalPrice}
		if r.Notes.Valid {
			line.Notes = r.Notes.String
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			line.CreatedAt = &t
		}
		if r.UpdatedAt.Valid {
			t := r.UpdatedAt.Time
			line.UpdatedAt = &t
		}
		out = append(out, line)
	}
	return c.JSON(http.StatusOK, out)
}

func CreatePurchaseOrderLine(c echo.Context) error {
	purchaseUID := strings.TrimSpace(c.Param("uid"))
	if purchaseUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid purchase order uid"})
	}
	var req models.PurchaseOrderLineCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if req.Quantity <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "quantity must be greater than zero"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "db not configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	orderRow, err := dbrepo.GetPurchaseOrderByUID(pool, purchaseUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "purchase order not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("purchase order lookup error: %v", err)})
	}

	var inventoryItemID *int
	if strings.TrimSpace(req.InventoryItemUID) != "" {
		inventoryRow, err := dbrepo.GetInventoryItemByUID(pool, strings.TrimSpace(req.InventoryItemUID))
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "inventory item uid not found"})
		}
		inventoryItemID = &inventoryRow.ID
	}

	uid, err := dbrepo.CreatePurchaseOrderLine(pool, req, orderRow.ID, inventoryItemID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": uid})
}

func UpdatePurchaseOrderLine(c echo.Context) error {
	uid := strings.TrimSpace(c.Param("line_uid"))
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid line uid"})
	}
	var req models.PurchaseOrderLineUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "db not configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	var inventoryItemID *int
	if req.InventoryItemUID != nil && strings.TrimSpace(*req.InventoryItemUID) != "" {
		inventoryRow, err := dbrepo.GetInventoryItemByUID(pool, strings.TrimSpace(*req.InventoryItemUID))
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "inventory item uid not found"})
		}
		inventoryItemID = &inventoryRow.ID
	}

	if err := dbrepo.UpdatePurchaseOrderLine(pool, uid, req, inventoryItemID); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

func DeletePurchaseOrderLine(c echo.Context) error {
	uid := strings.TrimSpace(c.Param("line_uid"))
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid line uid"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	if err := dbrepo.DeletePurchaseOrderLine(pool, uid); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete purchase order line: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}
