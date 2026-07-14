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

// ListVendorContacts godoc
// @Summary List Vendor Persons
// @Description Get all persons associated with a specific vendor (identified by vendor UID in path)
// @Tags VendorContacts
// @Produce json
// @Param vendor_uid path string true "Vendor UID"
// @Success 200 {array} models.VendorContact
// @Security BearerAuth
// @Router /vendors/{vendor_uid}/contacts [get]
func ListVendorContacts(c echo.Context) error {
	vuid := strings.Trim(c.Param("vendor_uid"), "/")
	if vuid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db not initailized"})
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

	rows, err := dbrepo.ListVendorContacts(pool, vendor.ID)
	if err != nil {
		c.Logger().Errorf("ListVendorPersons db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query vendor_persons: %v", err)})
	}
	var out []models.VendorContact
	for _, r := range rows {
		var v models.VendorContact
		v.UID = ""
		if r.UID.Valid {
			v.UID = r.UID.String
		}
		if r.Name.Valid {
			v.Name = r.Name.String
		}
		if r.Email.Valid {
			v.Email = r.Email.String
		}
		if r.Mobile.Valid {
			v.Mobile = r.Mobile.String
		}
		out = append(out, v)
	}
	return c.JSON(http.StatusOK, out)
}

// GetVendorContact godoc
// @Summary Get Vendor Contact
// @Description Get details for a specific vendor contact by vendor UID and contact UID
// @Tags VendorContacts
// @Produce json
// @Param vendor_uid path string true "Vendor UID"
// @Param contact_uid path string true "Contact UID"
// @Success 200 {object} models.VendorContact
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /vendors/{vendor_uid}/contacts/{contact_uid} [get]
func GetVendorContact(c echo.Context) error {
	vuid := strings.Trim(c.Param("vendor_uid"), "/")
	contactUID := strings.Trim(c.Param("contact_uid"), "/")
	if vuid == "" || contactUID == "" {
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
		c.Logger().Errorf("GetVendorContact vendor lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query vendor: %v", err)})
	}

	row, err := dbrepo.GetVendorContactByUID(pool, contactUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "contact not found"})
		}
		c.Logger().Errorf("GetVendorContact db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query contact: %v", err)})
	}

	var v models.VendorContact
	v.UID = ""
	if row.UID.Valid {
		v.UID = row.UID.String
	}
	if row.Name.Valid {
		v.Name = row.Name.String
	}
	if row.Email.Valid {
		v.Email = row.Email.String
	}
	if row.Mobile.Valid {
		v.Mobile = row.Mobile.String
	}

	return c.JSON(http.StatusOK, v)
}

// CreateVendorContact godoc
// @Summary Create Vendor Person
// @Description Create a new vendor person for a vendor (identified by vendor UID in path)
// @Tags VendorContacts
// @Accept json
// @Produce json
// @Param vid path string true "Vendor UID"
// @Param vendor_contact body models.VendorContactCreateRequest true "Vendor contact payload"
// @Security BearerAuth
// @Router /vendors/{vid}/contacts [post]
func CreateVendorContact(c echo.Context) error {
	vid := strings.Trim(c.Param("vendor_uid"), "/")
	if vid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusCreated, echo.Map{"note": "no db configured - vendor_contact payload received"})
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

	var req models.VendorContactCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	newID, err := dbrepo.CreateVendorContact(pool, req, vendor.ID)
	if err != nil {
		c.Logger().Errorf("CreateVendorContact db insert error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db insert vendor_contact: %v", err)})
	}
	return c.JSON(http.StatusCreated, echo.Map{"id": newID})
}

// UpdateVendorContact godoc
// @Summary Update Vendor Contact
// @Description Update a vendor contact by vendor UID and contact UID
// @Tags VendorContacts
// @Accept json
// @Produce json
// @Param vendor_uid path string true "Vendor UID"
// @Param contact_uid path string true "Contact UID"
// @Param vendor_contact body models.VendorContactUpdateRequest true "Vendor contact update payload"
// @Success 200 {object} map[string]string "{\"uid\":\"550e8400-e29b-41d4-a716-446655440000\"}"
// @Failure 400 {object} map[string]string "{\"error\":\"Invalid request\"}"
// @Failure 404 {object} map[string]string "{\"error\":\"Not found\"}"
// @Failure 500 {object} map[string]string "{\"error\":\"Internal server error\"}"
// @Security BearerAuth
// @Router /vendors/{vendor_uid}/contacts/{contact_uid} [patch]
func UpdateVendorContact(c echo.Context) error {
	vuid := strings.Trim(c.Param("vendor_uid"), "/")
	contactUID := strings.Trim(c.Param("contact_uid"), "/")
	if vuid == "" || contactUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	var req models.VendorContactUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would update vendor_contact", "uid": contactUID, "payload": req})
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

	// Verify contact exists
	contact, err := dbrepo.GetVendorContactByUID(pool, contactUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "contact uid not found"})
		}
		c.Logger().Errorf("GetVendorContactByUID db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query contact: %v", err)})
	}

	// Update the contact
	if err := dbrepo.UpdateVendorContact(pool, contact.ID, req); err != nil {
		c.Logger().Errorf("UpdateVendorContact db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db update vendor_contact: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": contactUID})
}

// DeleteVendorContact godoc
// @Summary Delete Vendor Person
// @Description Delete a vendor person by vendor UID and person UID
// @Tags VendorContacts
// @Param vid path string true "Vendor UID"
// @Param contact_uid path string true "Person UID"
// @Success 204 {string} string "No Content"
// @Security BearerAuth
// @Router /vendors/{vid}/contacts/{contact_uid} [delete]
func DeleteVendorContact(c echo.Context) error {
	vuid := strings.Trim(c.Param("vendor_uid"), "/")
	contact_uid := strings.Trim(c.Param("contact_uid"), "/")
	if vuid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.VendorContact{})
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

	contact, err := dbrepo.GetVendorContactByUID(pool, contact_uid)

	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "person Uid not found"})
		}
		c.Logger().Errorf("GetPerson db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query person: %v", err)})
	}

	if err := dbrepo.DeleteVendorContact(pool, contact.ID); err != nil {
		c.Logger().Errorf("DeleteVendorContact db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete vendor_contact: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}
