package controllers

import (
	"fmt"
	"net/http"
	"strings"

	dbrepo "app/db/repo"
	"app/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// ListVendorAddresses godoc
// @Summary List Vendor Addresses
// @Description Get all addresses associated with a specific vendor (identified by vendor UID in path)
// @Tags VendorAddresses
// @Produce json
// @Param vendor_uid path string true "Vendor UID"
// @Success 200 {array} models.VendorAddress
// @Security BearerAuth
// @Router /vendors/{vendor_uid}/addresses [get]
func ListVendorAddresses(c echo.Context) error {
	vuid := strings.Trim(c.Param("vendor_uid"), "/")
	if vuid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.VendorAddress{})
	}
	pool := dbi.(*pgxpool.Pool)

	vendor, err := dbrepo.GetVendorByUID(pool, vuid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "vendor id not found"})
		}
		c.Logger().Errorf("GetVendor db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query vendor: %v", err)})
	}

	rows, err := dbrepo.ListVendorAddresses(pool, vendor.ID)
	if err != nil {
		c.Logger().Errorf("ListVendorAddresses db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query vendor_addresses: %v", err)})
	}
	var out []models.VendorAddress
	for _, r := range rows {
		var v models.VendorAddress
		v.UID = ""
		if r.UID.Valid {
			v.UID = r.UID.String
		}
		v.IsPrimary = r.IsPrimary
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			v.CreatedAt = &t
		}
		if r.AddressLine1.Valid {
			v.AddressLine1 = r.AddressLine1.String
		}
		if r.AddressLine2.Valid {
			v.AddressLine2 = r.AddressLine2.String
		}
		if r.AddressLine3.Valid {
			v.AddressLine3 = r.AddressLine3.String
		}
		if r.Landmark.Valid {
			v.Landmark = r.Landmark.String
		}
		if r.City.Valid {
			v.City = r.City.String
		}
		if r.State.Valid {
			v.State = r.State.String
		}
		if r.Country.Valid {
			v.Country = r.Country.String
		}
		if r.PinCode.Valid {
			v.PinCode = r.PinCode.String
		}
		out = append(out, v)
	}
	return c.JSON(http.StatusOK, out)
}

// GetVendorAddress godoc
// @Summary Get Vendor Address
// @Description Get details for a specific vendor address by vendor UID and address UID
// @Tags VendorAddresses
// @Produce json
// @Param vendor_uid path string true "Vendor UID"
// @Param address_uid path string true "Address UID"
// @Success 200 {object} models.VendorAddress
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /vendors/{vendor_uid}/addresses/{address_uid} [get]
func GetVendorAddress(c echo.Context) error {
	vuid := strings.Trim(c.Param("vendor_uid"), "/")
	addressUID := strings.Trim(c.Param("address_uid"), "/")
	if vuid == "" || addressUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}
	pool := dbi.(*pgxpool.Pool)

	_, err := dbrepo.GetVendorByUID(pool, vuid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "vendor not found"})
		}
		c.Logger().Errorf("GetVendorAddress vendor lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query vendor: %v", err)})
	}

	row, err := dbrepo.GetVendorAddressByUID(pool, addressUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "address not found"})
		}
		c.Logger().Errorf("GetVendorAddress db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query address: %v", err)})
	}

	var v models.VendorAddress
	v.UID = ""
	if row.UID.Valid {
		v.UID = row.UID.String
	}
	v.IsPrimary = row.IsPrimary
	if row.CreatedAt.Valid {
		t := row.CreatedAt.Time
		v.CreatedAt = &t
	}
	if row.AddressLine1.Valid {
		v.AddressLine1 = row.AddressLine1.String
	}
	if row.AddressLine2.Valid {
		v.AddressLine2 = row.AddressLine2.String
	}
	if row.AddressLine3.Valid {
		v.AddressLine3 = row.AddressLine3.String
	}
	if row.Landmark.Valid {
		v.Landmark = row.Landmark.String
	}
	if row.City.Valid {
		v.City = row.City.String
	}
	if row.State.Valid {
		v.State = row.State.String
	}
	if row.Country.Valid {
		v.Country = row.Country.String
	}
	if row.PinCode.Valid {
		v.PinCode = row.PinCode.String
	}

	return c.JSON(http.StatusOK, v)
}

// // Get vendor address
// func GetVendorAddress(c echo.Context) error {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
// 	}
// 	dbi := c.Get("db")
// 	if dbi == nil {
// 		return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
// 	}
// 	db := dbi.(*sql.DB)
// 	r, err := dbrepo.GetVendorAddressByID(db, id)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
// 		}
// 		c.Logger().Errorf("GetVendorAddress db error: %v", err)
// 		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query vendor_address: %v", err)})
// 	}
// 	var v models.VendorAddress
// 	v.ID = strconv.Itoa(r.ID)
// 	if r.UID.Valid {
// 		v.UID = r.UID.String
// 	}
// 	v.VendorID = r.VendorID
// 	v.AddressID = r.AddressID
// 	v.IsPrimary = r.IsPrimary
// 	if r.CreatedAt.Valid {
// 		t := r.CreatedAt.Time
// 		v.CreatedAt = &t
// 	}
// 	return c.JSON(http.StatusOK, v)
// }

// CreateVendorAddress godoc
// @Summary Create Vendor Address
// @Description Create a new vendor address for a vendor (identified by vendor UID in path)
// @Tags VendorAddresses
// @Accept json
// @Produce json
// @Param vid path string true "Vendor UID"
// @Param vendor_address body models.VendorAddressCreateRequest true "Vendor address payload"
// @Security BearerAuth
// @Router /vendors/{vid}/addresses [post]
func CreateVendorAddress(c echo.Context) error {
	vid := strings.Trim(c.Param("vendor_uid"), "/")
	if vid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusCreated, echo.Map{"note": "no db configured - vendor_address payload received"})
	}
	pool := dbi.(*pgxpool.Pool)

	vendor, err := dbrepo.GetVendorByUID(pool, vid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "vendor id not found"})
		}
		c.Logger().Errorf("GetVendor db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query vendor: %v", err)})
	}

	var req models.VendorAddressCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	newID, err := dbrepo.CreateVendorAddress(pool, req, vendor.ID)
	if err != nil {
		c.Logger().Errorf("CreateVendorAddress db insert error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db insert vendor_address: %v", err)})
	}
	return c.JSON(http.StatusCreated, echo.Map{"id": newID})
}

// UpdateVendorAddress godoc
// @Summary Update Vendor Address
// @Description Update a vendor address by vendor UID and address UID
// @Tags VendorAddresses
// @Accept json
// @Produce json
// @Param vendor_uid path string true "Vendor UID"
// @Param address_uid path string true "Address UID"
// @Param vendor_address body models.VendorAddressUpdateRequest true "Vendor address update payload"
// @Success 200 {object} map[string]string "{\"uid\":\"550e8400-e29b-41d4-a716-446655440000\"}"
// @Failure 400 {object} map[string]string "{\"error\":\"Invalid request\"}"
// @Failure 404 {object} map[string]string "{\"error\":\"Not found\"}"
// @Failure 500 {object} map[string]string "{\"error\":\"Internal server error\"}"
// @Security BearerAuth
// @Router /vendors/{vendor_uid}/addresses/{address_uid} [patch]
func UpdateVendorAddress(c echo.Context) error {
	vuid := strings.Trim(c.Param("vendor_uid"), "/")
	addressUID := strings.Trim(c.Param("address_uid"), "/")
	if vuid == "" || addressUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	var req models.VendorAddressUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would update vendor_address", "uid": addressUID, "payload": req})
	}
	pool := dbi.(*pgxpool.Pool)

	// Verify vendor exists
	_, err := dbrepo.GetVendorByUID(pool, vuid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "vendor id not found"})
		}
		c.Logger().Errorf("GetVendor db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query vendor: %v", err)})
	}

	// Verify address exists
	address, err := dbrepo.GetVendorAddressByUID(pool, addressUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "address uid not found"})
		}
		c.Logger().Errorf("GetVendorAddressByUID db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query address: %v", err)})
	}

	// Update the address
	if err := dbrepo.UpdateVendorAddress(pool, address.ID, req); err != nil {
		c.Logger().Errorf("UpdateVendorAddress db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db update vendor_address: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": addressUID})
}

// Delete vendor address
// DeleteVendorAddress godoc
// @Summary Delete Vendor Address
// @Description Delete a vendor address by vendor UID and address UID
// @Tags VendorAddresses
// @Param vid path string true "Vendor UID"
// @Param address_uid path string true "Address UID"
// @Success 204 {string} string "No Content"
// @Security BearerAuth
// @Router /vendors/{vid}/addresses/{address_uid} [delete]
func DeleteVendorAddress(c echo.Context) error {
	vuid := strings.Trim(c.Param("vendor_uid"), "/")
	address_uid := strings.Trim(c.Param("address_uid"), "/")
	if vuid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.VendorAddress{})
	}
	pool := dbi.(*pgxpool.Pool)

	_, err := dbrepo.GetVendorByUID(pool, vuid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "vendor id not found"})
		}
		c.Logger().Errorf("GetVendor db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query vendor: %v", err)})
	}

	address, err := dbrepo.GetVendorAddressByUID(pool, address_uid)

	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "address Uid not found"})
		}
		c.Logger().Errorf("Getaddress db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query address: %v", err)})
	}

	if err := dbrepo.DeleteVendorAddress(pool, address.ID); err != nil {
		c.Logger().Errorf("DeleteVendorAddress db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete vendor_address: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}
