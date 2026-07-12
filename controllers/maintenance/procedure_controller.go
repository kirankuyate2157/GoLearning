package controllers

import (
	dbrepo "app/db/repo/maintenance"
	models "app/models/maintenance"
	"fmt"
	"net/http"
	"strings"

	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// ListTemplateProcedures godoc
// @Summary Get Template Procedures
// @Description Get list of template procedures
// @Tags Maintenance Template Procedures
// @Param uid path string true "Maintenance Template UID"
// @Produce json
// @Success 200 {array} models.TemplateProcedure
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /maintenance/templates/{uid}/procedures [get]
func ListTemplateProcedures(c echo.Context) error {
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
		c.Logger().Error("ListTemplateProcedures: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	queryParams := models.MaintenanceTemplateParams{UID: uid}

	maintenanceRows, _, err := dbrepo.GetMaintenanceTemplate(pool, queryParams, false)
	if err != nil || len(maintenanceRows) == 0 {
		if err == pgx.ErrNoRows || len(maintenanceRows) == 0 {
			return c.JSON(http.StatusNotFound, echo.Map{"error": fmt.Sprintf("UID: %v not found", uid)})
		}
		c.Logger().Errorf("DeleteMaintenanceTemplate lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query template: %v", err)})
	}

	rows, err := dbrepo.GetTemplateProcedure(pool, dbrepo.TemplateProcedureQP{TemplateID: &maintenanceRows[0].ID})
	if err != nil {
		c.Logger().Errorf("ListTemplateProcedures db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query procedures: %v", err)})
	}
	var out []models.TemplateProcedure

	for _, r := range rows {
		out = append(out, models.TemplateProcedure{
			UID:         r.UID,
			Name:        r.Name,
			Description: r.Description,
			TemplateUID: r.TemplateUID,
		})
	}
	if out == nil {
		out = []models.TemplateProcedure{}
	}
	return c.JSON(http.StatusOK, echo.Map{
		"total_templates": len(rows),
		"procedures":      out,
	})

}

// GetTemplateProcedure godoc
// @Summary Get Maintenance Template Procedure
// @Description Get template procedure by UID
// @Tags Maintenance Template Procedures
// @Produce json
// @Param uid path string true "Procedure UID"
// @Success 200 {object} models.TemplateProcedure
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance/templates/procedures/{uid} [get]
func GetTemplateProcedure(c echo.Context) error {
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
		c.Logger().Error("GetTemplateProcedure: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	queryParams := dbrepo.TemplateProcedureQP{ProcedureUID: &uid}

	// Get procedure row(s) (returns slice)
	rows, err := dbrepo.GetTemplateProcedure(pool, queryParams)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("GetTemplateProcedure db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("db query procedure: %v", err),
		})
	}
	if len(rows) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
	}
	r := rows[0]
	// Convert template_id from string → int
	tid, err := strconv.Atoi(r.TemplateID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid template_id, must be a number",
		})
	}
	templatequeryParams := models.MaintenanceTemplateParams{}
	templatequeryParams.ID = tid

	//  Get template row using integer ID
	rows1, _, err := dbrepo.GetMaintenanceTemplate(pool, templatequeryParams, false)
	templateRow := rows1[0]
	var templateUID string
	if err == nil {
		templateUID = templateRow.UID
	} else {
		templateUID = ""
	}
	//  Build final response
	f := models.TemplateProcedure{
		UID:         r.UID,
		Name:        r.Name,
		Description: r.Description,
		TemplateUID: templateUID,
	}
	return c.JSON(http.StatusOK, f)
}

// CreateTemplateProcedure godoc
// @Summary Create Maintenance Template Procedure
// @Description Create a new template procedure
// @Tags Maintenance Template Procedures
// @Accept json
// @Produce json
// @Param uid path string true "Maintenance Template UID"
// @Param procedure body models.TemplateProcedureCreateRequest true "Procedure payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance/templates/{uid}/procedures [post]
func CreateTemplateProcedure(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var req models.TemplateProcedureCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - procedure payload received", "procedure": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("CreateTemplateProcedure: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	templateUID := uid
	queryParams := models.MaintenanceTemplateParams{}
	queryParams.UID = templateUID

	rows, _, err := dbrepo.GetMaintenanceTemplate(pool, queryParams, false)
	if err != nil || len(rows) == 0 {
		c.Logger().Errorf("CreateTemplateProcedure: template lookup error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid template UID"})
	}
	templateRow := rows[0]

	// Use integer ID for DB
	newUID, err := dbrepo.CreateTemplateProcedure(pool, dbrepo.TemplateProcedureInsert{
		Name:        req.Name,
		Description: req.Description,
		TemplateID:  templateRow.ID,
	})
	if err != nil {
		c.Logger().Errorf("CreateTemplateProcedure db insert error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db insert procedure: %v", err)})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

// UpdateTemplateProcedure godoc
// @Summary Update Maintenance Template Procedure
// @Description Update template procedure by UID
// @Tags Maintenance Template Procedures
// @Accept json
// @Produce json
// @Param uid path string true "Procedure UID"
// @Param procedure body models.TemplateProcedureUpdateRequest true "Procedure update payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance/templates/procedures/{uid} [patch]
func UpdateTemplateProcedure(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var req models.TemplateProcedureUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would update procedure", "uid": uid, "payload": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("UpdateTemplateProcedure: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	reqInt := req
	queryParams := models.MaintenanceTemplateParams{}
	queryParams.UID = *req.TemplateUID
	fmt.Println(queryParams.UID)
	if req.TemplateUID != nil {
		rows, _, err := dbrepo.GetMaintenanceTemplate(pool, queryParams, false)
		templateRow := rows[0]

		if err != nil {
			c.Logger().Errorf("UpdateTemplateProcedure: template lookup error: %v", err)
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid template UID"})
		}
		intID := fmt.Sprintf("%d", templateRow.ID)
		reqInt.TemplateUID = &intID
	}
	if err := dbrepo.UpdateTemplateProcedure(pool, uid, reqInt); err != nil {
		c.Logger().Errorf("UpdateTemplateProcedure db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db update procedure: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

// DeleteTemplateProcedure godoc
// @Summary Delete Maintenance Template Procedure
// @Description Delete a template procedure by UID
// @Tags Maintenance Template Procedures
// @Produce json
// @Param uid path string true "Procedure UID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance/templates/procedures/{uid} [delete]
func DeleteTemplateProcedure(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would not delete procedure", "uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteTemplateProcedure: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	queryParams := dbrepo.TemplateProcedureQP{ProcedureUID: &uid}

	if _, err := dbrepo.GetTemplateProcedure(pool, queryParams); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("DeleteTemplateProcedure lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query procedure: %v", err)})
	}
	if err := dbrepo.DeleteTemplateProcedure(pool, uid); err != nil {
		c.Logger().Errorf("DeleteTemplateProcedure db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete procedure: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Procedure deleted successfully"})
}
