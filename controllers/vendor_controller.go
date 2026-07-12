package controllers

import (
	"fmt"
	"net/http"
	"strings"

	dbrepo "app/db/repo"
	dbreposettings "app/db/repo/settings"
	"app/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// GetVendors godoc
// @Summary Get Vendors
// @Description Get list of vendors (JWT protected)
// @Tags Vendors
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Vendor
// @Router /vendors [get]
// @Param vendor_type query string false "Filter by vendor type name"
func GetVendors(c echo.Context) error {
	typeFilter := strings.TrimSpace(c.QueryParam("vendor_type"))

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.Vendor{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	rows, err := dbrepo.ListVendors(pool, typeFilter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query vendors: %v", err)})
	}

	var out []models.Vendor
	for _, r := range rows {
		var v models.Vendor
		if r.UID.Valid {
			v.UID = r.UID.String
		}
		v.Name = r.Name
		if r.Website.Valid {
			v.Website = r.Website.String
		}
		if r.Email.Valid {
			v.Email = r.Email.String
		}
		if r.Phone.Valid {
			v.PhoneNo = r.Phone.String
		}
		if r.JoinDate.Valid {
			t := r.JoinDate.Time
			v.JoinDate = &t
		}
		if r.GSTNo.Valid {
			v.GSTNo = r.GSTNo.String
		}
		if r.GSTType.Valid {
			v.GSTTypeUID = r.GSTType.String
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			v.CreatedAt = &t
		}

		// Fetch files for this vendor
		entityType := "vendor"
		entityId := int64(r.ID)
		fileFilters := dbrepo.FileListFilters{EntityType: &entityType, EntityID: &entityId}
		fileRows, _, _, _ := dbrepo.ListFilesWithFilters(pool, fileFilters)
		var files []models.File
		for _, file := range fileRows {
			files = append(files, models.File{
				UID:          file.UID,
				OriginalName: file.OriginalName,
				Extension:    file.Extension,
				MimeType:     file.MimeType,
				SizeBytes:    file.SizeBytes,
				UploadedAt:   file.UploadedAt,
				UsageType:    file.UsageType,
				IsPrimary:    file.IsPrimary,
			})
		}
		v.Files = files

		out = append(out, v)
	}

	return c.JSON(http.StatusOK, out)
}

// GetVendor godoc
// @Summary Get Vendor
// @Description Get vendor by UID (JWT protected)
// @Tags Vendors
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Vendor UID (UID)"
// @Success 200 {object} models.Vendor
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /vendors/{uid} [get]
func GetVendor(c echo.Context) error {
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
		c.Logger().Error("GetVendor: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	r, err := dbrepo.GetVendorByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("GetVendor db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query vendor: %v", err)})
	}
	var v models.Vendor
	if r.UID.Valid {
		v.UID = r.UID.String
	}
	v.Name = r.Name
	if r.Website.Valid {
		v.Website = r.Website.String
	}
	if r.Email.Valid {
		v.Email = r.Email.String
	}
	if r.Phone.Valid {
		v.PhoneNo = r.Phone.String
	}
	if r.JoinDate.Valid {
		t := r.JoinDate.Time
		v.JoinDate = &t
	}
	if r.GSTNo.Valid {
		v.GSTNo = r.GSTNo.String
	}
	if r.GSTType.Valid {
		v.GSTTypeUID = r.GSTType.String
	}
	if r.CreatedAt.Valid {
		t := r.CreatedAt.Time
		v.CreatedAt = &t
	}

	// Fetch files for this vendor
	entityType := "vendor"
	entityId := int64(r.ID)
	fileFilters := dbrepo.FileListFilters{EntityType: &entityType, EntityID: &entityId}
	fileRows, _, _, _ := dbrepo.ListFilesWithFilters(pool, fileFilters)
	var files []models.File
	for _, file := range fileRows {
		files = append(files, models.File{
			UID:          file.UID,
			OriginalName: file.OriginalName,
			Extension:    file.Extension,
			MimeType:     file.MimeType,
			SizeBytes:    file.SizeBytes,
			UploadedAt:   file.UploadedAt,
			UsageType:    file.UsageType,
			IsPrimary:    file.IsPrimary,
		})
	}
	v.Files = files

	return c.JSON(http.StatusOK, v)
}

// CreateVendor godoc
// @Summary Create Vendor
// @Description Create a new vendor (JWT protected)
// @Tags Vendors
// @Accept json
// @Produce json
// @Param vendor body models.VendorCreateRequest true "Vendor payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /vendors [post]
func CreateVendor(c echo.Context) error {
	var req models.VendorCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - vendor payload received", "vendor": req})
	}
	pool := dbi.(*pgxpool.Pool)

	// Get tax ID from tax UID
	var taxID int
	if req.GSTTypeUID != "" {
		taxRow, err := dbreposettings.GetTaxByUID(pool, req.GSTTypeUID)
		if err != nil {
			c.Logger().Errorf("CreateVendor GetTaxByUID error: %v", err)
			return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("tax not found: %v", err)})
		}
		taxID = int(taxRow.ID)
	}

	newID, err := dbrepo.CreateVendor(pool, req, taxID)
	if err != nil {
		c.Logger().Errorf("CreateVendor db insert error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db insert vendor: %v", err)})
	}
	return c.JSON(http.StatusCreated, echo.Map{"id": newID})
}

// UpdateVendor godoc
// @Summary Update Vendor
// @Description Update an existing vendor by id (JWT protected)
// @Tags Vendors
// @Accept json
// @Produce json
// @Param id path int true "Vendor ID"
// @Param vendor body models.VendorUpdateRequest true "Vendor update payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /vendors/{id} [patch]
func UpdateVendor(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var req models.VendorUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would update vendor", "id": uid, "payload": req})
	}
	pool := dbi.(*pgxpool.Pool)

	// Get tax ID from tax UID if provided
	var taxID *int
	if req.GSTTypeUID != nil && *req.GSTTypeUID != "" {
		taxRow, err := dbreposettings.GetTaxByUID(pool, *req.GSTTypeUID)
		if err != nil {
			c.Logger().Errorf("UpdateVendor GetTaxByUID error: %v", err)
			return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("tax not found: %v", err)})
		}
		id := int(taxRow.ID)
		taxID = &id
	}

	if err := dbrepo.UpdateVendor(pool, uid, req, taxID); err != nil {
		c.Logger().Errorf("UpdateVendor db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db update vendor: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"id": uid})
}

// DeleteVendor godoc
// @Summary Delete Vendor
// @Description Delete a vendor by id (JWT protected)
// @Tags Vendors
// @Produce json
// @Param id path int true "Vendor ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /vendors/{id} [delete]
func DeleteVendor(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would delete vendor", "id": uid})
	}
	pool := dbi.(*pgxpool.Pool)
	if err := dbrepo.DeleteVendor(pool, uid); err != nil {
		c.Logger().Errorf("DeleteVendor db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete vendor: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}
