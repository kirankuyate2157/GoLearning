package controllers

import (
	"app/db/repo"
	"app/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// CreateVendorType godoc
// @Summary Assign Vendor Type
// @Description Assign a vendor type to a specific vendor using vendor_uid
// @Tags Vendor Types
// @Accept json
// @Produce json
// @Param vendor_uid path string true "Vendor UID"
// @Param payload body models.VendorTypeRequest true "Vendor Type Request"
// @Success 201 {object} models.VendorType
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /vendors/{vendor_uid}/vendorsTypes [post]
func CreateVendorType(c echo.Context) error {
	vendorUID := c.Param("vendor_uid")

	var req models.VendorTypeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON body"})
	}

	if req.VendorTypeUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "vendor_type_uid is required"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("CreateVendorType: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	vendorType, err := repo.CreateVendorType(pool, vendorUID, req.VendorTypeUID)
	if err != nil {
		c.Logger().Errorf("CreateVendorType db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db insert vendor type: %v", err)})
	}

	return c.JSON(http.StatusCreated, vendorType)
}

// GetVendorTypes godoc
// @Summary Get Vendor Types
// @Description Retrieve all vendor types assigned to a vendor
// @Tags Vendor Types
// @Produce json
// @Param vendor_uid path string true "Vendor UID"
// @Success 200 {array} string
// @Failure 500 {object} map[string]string
// @Router /vendors/{vendor_uid}/vendorsTypes [get]
func GetVendorTypes(c echo.Context) error {
	vuid := strings.TrimSpace(c.Param("vendor_uid"))
	if vuid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid vendor_uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("GetVendorTypes: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	rows, err := repo.GetVendorTypes(pool, vuid)
	if err != nil {
		c.Logger().Errorf("GetVendorTypes db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query vendor types: %v", err)})
	}

	var out []models.VendorType

	for _, r := range rows {
		var vt models.VendorType

		if r.TypeUID.Valid {
			vt.TypeUID = r.TypeUID.String
		}
		if r.VendorAssignUID.Valid {
			vt.VendorAssignUID = r.VendorAssignUID.String
		}
		if r.Name.Valid {
			vt.Name = r.Name.String
		}
		if r.Description.Valid {
			vt.Description = r.Description.String
		}
		if r.Category.Valid {
			vt.Category = r.Category.String
		}

		out = append(out, vt)
	}

	if out == nil {
		out = []models.VendorType{}
	}

	return c.JSON(http.StatusOK, out)
}

// DeleteVendorType godoc
// @Summary Delete Vendor Type
// @Description Remove a specific vendor type assignment by vendorAssign_uid
// @Tags Vendor Types
// @Param vendor_uid path string true "Vendor UID"
// @Param vendorAssign_uid path string true "Vendor Assign UID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /vendors/{vendor_uid}/vendorsTypes/{vendorAssign_uid} [delete]
func DeleteVendorType(c echo.Context) error {
	vendorUID := c.Param("vendor_uid")
	vendorAssignUID := c.Param("vendorAssign_uid")

	if vendorUID == "" || vendorAssignUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "vendor_uid and vendorAssign_uid are required"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteVendorType: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	if err := repo.DeleteVendorType(pool, vendorUID, vendorAssignUID); err != nil {
		c.Logger().Errorf("DeleteVendorType db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete vendor type: %v", err)})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Vendor type deleted successfully"})
}
