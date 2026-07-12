package controllers

import (
	dbrepo "app/db/repo/maintenance"
	models "app/models/maintenance"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// ListMaintenanceTypes godoc
// @Summary Get MaintenanceTypes
// @Description Get list of Maintenance_types
// @Tags MaintenanceTypes
// @Produce json
// @Success 200 {array} models.MaintenanceType
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /maintenance/type [get]
func ListMaintenanceTypes(c echo.Context) error {
	// Get DB from context
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "database not found in context",
		})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("ListMaintenanceTypes: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "invalid db object in context",
		})
	}
	queryParams := models.MaintenanceType{}

	// Fetch rows from DB
	rows, err := dbrepo.GetMaintenanceType(pool, queryParams)
	if err != nil {
		c.Logger().Errorf("ListMaintenanceTypes db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("db query Maintenance types: %v", err),
		})
	}

	// Final output
	var output []models.MaintenanceType

	// Loop through DB rows
	for _, r := range rows {
		var ft models.MaintenanceType

		// UID
		if r.UID.Valid {
			ft.UID = r.UID.String
		}

		// Name
		if r.Name.Valid {
			ft.Name = r.Name.String
		}

		// Description
		if r.Description.Valid {
			ft.Description = r.Description.String
		}

		// -------------------------------------------------

		output = append(output, ft)
	}

	// Always return array, never null
	if output == nil {
		output = []models.MaintenanceType{}
	}

	return c.JSON(http.StatusOK, output)
}

// GetMaintenanceType godoc
// @Summary Get MaintenanceType
// @Description Get a Maintenance type by UID
// @Tags MaintenanceTypes
// @Produce json
// @Param uid path string true "MaintenanceType UID"
// @Success 200 {object} models.MaintenanceType
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance/type/{uid} [get]
func GetMaintenanceType(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("GetMaintenanceType: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	queryParams := models.MaintenanceType{}
	queryParams.UID = uid

	fmt.Println(queryParams)
	// Fetch record from DB
	rows, err := dbrepo.GetMaintenanceType(pool, queryParams)
	r := rows[0]
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("GetMaintenanceType db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query Maintenance types: %v", err)})
	}

	// Convert DB row to model
	var ft models.MaintenanceType

	if r.UID.Valid {
		ft.UID = r.UID.String
	}
	if r.Name.Valid {
		ft.Name = r.Name.String
	}
	if r.Description.Valid {
		ft.Description = r.Description.String
	}

	// -------------------------------------------------

	return c.JSON(http.StatusOK, ft)
}

// CreateMaintenanceType godoc
// @Summary Create MaintenanceType
// @Description Create a new Maintenance type
// @Tags MaintenanceTypes
// @Accept json
// @Produce json
// @Param Maintenance_type body models.MaintenanceTypeCreateRequest true "MaintenanceType payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance/type [post]
func CreateMaintenanceType(c echo.Context) error {
	var req models.MaintenanceTypeCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - Maintenance type payload received", "Maintenance_type": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("CreateMaintenanceType: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	newUID, err := dbrepo.CreateMaintenanceType(pool, req)
	if err != nil {
		c.Logger().Errorf("CreateMaintenanceType db insert error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db insert Maintenance type: %v", err)})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

// UpdateMaintenanceType godoc
// @Summary Update MaintenanceType
// @Description Update an existing Maintenance type by UID
// @Tags MaintenanceTypes
// @Accept json
// @Produce json
// @Param uid path string true "MaintenanceType UID"
// @Param Maintenance_type body models.MaintenanceTypeUpdateRequest true "MaintenanceType update payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance/type/{uid} [patch]
func UpdateMaintenanceType(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var req models.MaintenanceTypeUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would update Maintenance type", "uid": uid, "payload": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("UpdateMaintenanceType: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	if err := dbrepo.UpdateMaintenanceType(pool, uid, req); err != nil {
		c.Logger().Errorf("UpdateMaintenanceType db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db update Maintenance type: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

// DeleteMaintenanceType godoc
// @Summary Delete MaintenanceType
// @Description Delete a Maintenance type by UID
// @Tags MaintenanceTypes
// @Produce json
// @Param uid path string true "MaintenanceType UID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance/type/{uid} [delete]
func DeleteMaintenanceType(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would not delete Maintenance type", "uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteMaintenanceType: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	queryParams := models.MaintenanceType{}
	queryParams.UID = uid

	if _, err := dbrepo.GetMaintenanceType(pool, queryParams); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("DeleteMaintenanceType lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query Maintenance type: %v", err)})
	}

	if err := dbrepo.DeleteMaintenanceType(pool, uid); err != nil {
		c.Logger().Errorf("DeleteMaintenanceType db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete Maintenance type: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Maintenance type deleted successfully"})
}
