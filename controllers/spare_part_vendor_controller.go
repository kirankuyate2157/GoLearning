package controllers

import (
	"net/http"
	"strings"

	dbrepo "app/db/repo"
	"app/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// GetSparePartVendorList godoc
// @Summary Get vendors for a spare part
// @Tags SparePartVendor
// @Produce json
// @Param sparepart_uid path string true "Spare Part UID"
// @Success 200 {array} repo.SparePartVendorRow
// @Router /spare-parts/{sparepart_uid}/vendors [get]
func GetSparePartVendorList(c echo.Context) error {
	spareUID := c.Param("sparepart_uid")

	if spareUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid sparepart_uid"})
	}

	pool := c.Get("db").(*pgxpool.Pool)

	rows, err := dbrepo.ListSparePartVendors(pool, spareUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if rows == nil {
		rows = []dbrepo.SparePartVendorRow{}
	}

	return c.JSON(http.StatusOK, rows)
}

// CreateSparePartVendor godoc
// @Summary Add vendor to spare part
// @Tags SparePartVendor
// @Accept json
// @Produce json
// @Param sparepart_uid path string true "Spare Part UID"
// @Param body body models.SparePartVendorCreateRequest true "Vendor assignment payload"
// @Success 201 {object} map[string]string
// @Router /spare-parts/{sparepart_uid}/vendors [post]
func CreateSparePartVendor(c echo.Context) error {
	spareUID := c.Param("sparepart_uid")
	if spareUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid sparepart_uid"})
	}

	var req models.SparePartVendorCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if strings.TrimSpace(req.VendorUID) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "vendor_uid is required"})
	}

	pool := c.Get("db").(*pgxpool.Pool)

	// resolve spare_part_id
	sp, err := dbrepo.GetSparePartByUID(pool, spareUID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "spare part not found"})
	}

	// resolve vendor_id
	v, err := dbrepo.GetVendorByUID(pool, req.VendorUID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "vendor not found"})
	}

	modelNo := ""
	if req.ModelNo != nil {
		modelNo = *req.ModelNo
	}

	newUID, err := dbrepo.CreateSparePartVendor(pool, sp.ID, v.ID, modelNo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

// UpdateSparePartVendor godoc
// @Summary Update spare_part_vendor entry
// @Tags SparePartVendor
// @Accept json
// @Produce json
// @Param sparepart_uid path string true "Spare Part UID"
// @Param sparepartvendor_uid path string true "SparePartVendor UID"
// @Param body body models.SparePartVendorUpdateRequest true "Update payload"
// @Success 200 {object} map[string]interface{}
// @Router /spare-parts/{sparepart_uid}/vendors/{sparepartvendor_uid} [patch]
func UpdateSparePartVendor(c echo.Context) error {
	sparePartUID := c.Param("sparepart_uid")
	vendorUID := c.Param("sparepartvendor_uid")

	if sparePartUID == "" || vendorUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing parameters"})
	}

	var req models.SparePartVendorUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if req.ModelNo == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "model_no required"})
	}

	pool := c.Get("db").(*pgxpool.Pool)

	sparePart, err := dbrepo.GetSparePartByUID(pool, sparePartUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	vendor, err := dbrepo.GetVendorByUID(pool, vendorUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	affected, err := dbrepo.UpdateSparePartVendor(pool, sparePart.ID, vendor.ID, *req.ModelNo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if affected == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "mapping not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"sparepart_uid": sparePartUID,
	})
}

// DeleteSparePartVendor godoc
// @Summary Delete spare_part_vendor entry
// @Tags SparePartVendor
// @Param sparepart_uid path string true "Spare Part UID"
// @Param sparepartvendor_uid path string true "SparePartVendor UID"
// @Success 204
// @Router /spare-parts/{sparepart_uid}/vendors/{sparepartvendor_uid} [delete]
func DeleteSparePartVendor(c echo.Context) error {
	vendorUID := c.Param("sparepartvendor_uid")
	sparePartUID := c.Param("sparepart_uid")

	if vendorUID == "" || sparePartUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing parameters"})
	}

	pool := c.Get("db").(*pgxpool.Pool)

	vendor, err := dbrepo.GetVendorByUID(pool, vendorUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	sparePart, err := dbrepo.GetSparePartByUID(pool, sparePartUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	affected, err := dbrepo.DeleteSparePartVendor(pool, sparePart.ID, vendor.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if affected == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "mapping not found"})
	}

	return c.NoContent(http.StatusNoContent)
}
