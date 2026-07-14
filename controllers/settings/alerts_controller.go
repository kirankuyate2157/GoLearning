package controllers

import (
	dbrepo "app/db/repo/settings"
	"app/models/settings"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// ListAlertsTypes godoc
// @Summary Get Alerts Types List
// @Description Get Alerts Types List
// @Tags Alerts
// @Produce json
// @Success 200 {object} models.AlertsTypesGet
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /settings/alerts/types [get]
func ListAlertsTypes(c echo.Context) error {
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "database not found in context",
		})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("ListAlertsTypes: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "invalid db object in context",
		})
	}

	rows, err := dbrepo.ListAlertsTypes(pool, dbrepo.AlertsTypeFilters{})
	if err != nil {
		c.Logger().Errorf("ListAlertsTypes db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("db query alert types: %v", err),
		})
	}

	var output []models.AlertsTypesGet

	// Loop through DB rows
	for _, r := range rows {
		var ft models.AlertsTypesGet

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

		// -------- JSON UNMARSHAL FOR STRUCTURE --------
		if r.Structure.Valid && r.Structure.String != "" {
			var jsonMap []map[string]interface{}

			if err := json.Unmarshal([]byte(r.Structure.String), &jsonMap); err != nil {
				c.Logger().Errorf("failed to unmarshal structure JSON for UID %s: %v", ft.UID, err)
				ft.Structure = []map[string]interface{}{} // fallback
			} else {
				ft.Structure = jsonMap
			}

		} else {
			ft.Structure = []map[string]interface{}{}
		}
		// -------------------------------------------------

		output = append(output, ft)
	}

	// Always return array, never null
	if output == nil {
		output = []models.AlertsTypesGet{}
	}

	return c.JSON(http.StatusOK, output)
}

// ListAlerts godoc
// @Summary Get  Alerts List
// @Description Get alert by UID
// @Tags Alerts
// @Produce json
// @Success 200 {object} models.AlertsGet
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /settings/alerts [get]
func ListAlerts(c echo.Context) error {
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "database not found in context",
		})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("ListAlerts: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "invalid db object in context",
		})
	}
	queryParams := dbrepo.AlertsFilters{}

	rows, err := dbrepo.ListAlerts(pool, queryParams)
	if err != nil {
		c.Logger().Errorf("ListAlerts db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("db query alerts: %v", err),
		})
	}

	var output []models.AlertsGet

	// Loop through DB rows
	for _, r := range rows {
		var ft models.AlertsGet

		if r.UID.Valid {
			ft.UID = r.UID.String
		}

		if r.Name.Valid {
			ft.Name = r.Name.String
		}

		if r.TypeUID.Valid {
			ft.TypeUID = r.TypeUID.String
		}

		// -------- JSON UNMARSHAL FOR STRUCTURE --------
		if r.Config.Valid && r.Config.String != "" {
			var jsonMap []map[string]interface{}

			if err := json.Unmarshal([]byte(r.Config.String), &jsonMap); err != nil {
				c.Logger().Errorf("failed to unmarshal structure JSON for UID %s: %v", ft.UID, err)
				ft.Config = []map[string]interface{}{} // fallback
			} else {
				ft.Config = jsonMap
			}

		} else {
			ft.Config = []map[string]interface{}{}
		}
		// -------------------------------------------------

		output = append(output, ft)
	}

	// Always return array, never null
	if output == nil {
		output = []models.AlertsGet{}
	}

	return c.JSON(http.StatusOK, output)
}

// CreateAlert godoc
// @Summary Create Alert
// @Description Create a new  alert
// @Tags Alerts
// @Accept json
// @Produce json
// @Param feed_alert body models.AlertsPost true "FeedAlert payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /settings/alerts [post]
func CreateAlerts(c echo.Context) error {
	var req models.AlertsPost
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - feed alert payload received", "alerts": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("CreateAlert: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	alertTypeQP := dbrepo.AlertsTypeFilters{}
	alertTypeQP.UID = &req.TypeUID

	alertTypesList, err := dbrepo.ListAlertsTypes(pool, alertTypeQP)
	if err != nil {
		c.Logger().Errorf("CreateAlerts ListAlertsTypes db error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db query alert types: %v", err)})
	}

	if len(alertTypesList) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("alert type with UID %s not found", req.TypeUID)})
	}

	data := dbrepo.AlertsCreate{
		Name:      req.Name,
		TypeID:    alertTypesList[0].ID,
		Structure: req.Structure,
	}

	newUID, err := dbrepo.CreateAlerts(pool, data)
	if err != nil {
		c.Logger().Errorf("CreateAlerts db insert error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db insert alerts: %v", err)})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})

}

// UpdateAlert godoc
// @Summary Update Alert
// @Description Update alert by UID
// @Tags Alerts
// @Accept json
// @Produce json
// @Param uid path string true "Alert UID"
// @Param feed_alert body models.AlertsPatch true "Alert update payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /settings/alerts/{uid} [patch]
func UpdateAlert(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}
	var req models.AlertsPatch
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would update feed alert", "uid": uid, "payload": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("UpdateAlert: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	// Check alert exists
	queryParams := dbrepo.AlertsFilters{}
	queryParams.UID = &uid

	alertList, err := dbrepo.ListAlerts(pool, queryParams)
	if err != nil {
		c.Logger().Errorf("ListAlerts db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db query alerts error"})
	}

	if len(alertList) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Alert not found"})
	}

	alertID := alertList[0].ID

	// Get TypeID
	alertTypeQP := dbrepo.AlertsTypeFilters{}
	alertTypeQP.UID = req.TypeUID

	alertTypesList, err := dbrepo.ListAlertsTypes(pool, alertTypeQP)
	if err != nil {
		c.Logger().Errorf("CreateAlerts ListAlertsTypes db error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db query alert types: %v", err)})
	}

	if len(alertTypesList) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("alert type with UID %s not found", *req.TypeUID)})
	}

	data := dbrepo.AlertsUpdate{
		Name:      req.Name,
		TypeID:    &alertTypesList[0].ID,
		Structure: req.Structure,
	}

	if err := dbrepo.UpdateAlert(pool, alertID, data); err != nil {
		c.Logger().Errorf("UpdateAlert db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db update feed alert: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})

}

// DeleteAlert godoc
// @Summary Delete Alert
// @Description Delete a alert by UID
// @Tags Alerts
// @Produce json
// @Param uid path string true "Alert UID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /settings/alerts/{uid} [delete]
func DeleteAlert(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")

	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid alert uid"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would not delete  alert", "uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteAlert: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	// Check alert exists
	queryParams := dbrepo.AlertsFilters{}
	queryParams.UID = &uid

	alertList, err := dbrepo.ListAlerts(pool, queryParams)
	if err != nil {
		c.Logger().Errorf("ListAlerts db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db query alerts error"})
	}

	if len(alertList) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Alert not found"})
	}

	id := alertList[0].ID

	if err := dbrepo.DeleteAlert(pool, id); err != nil {
		c.Logger().Errorf("DeleteAlert db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete feed alert: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "alert deleted successfully"})
}
