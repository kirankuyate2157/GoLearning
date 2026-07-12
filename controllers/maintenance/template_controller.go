package controllers

import (
	typeRepo "app/db/repo"
	dbrepo "app/db/repo/maintenance"
	models "app/models/maintenance"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// ListMaintenanceTemplates godoc
// @Summary Get Maintenance Templates
// @Description Get list of maintenance templates
// @Tags MaintenanceTemplates
// @Produce json
// @Success 200 {array} models.MaintenanceTemplate
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /maintenance/templates [get]
func ListMaintenanceTemplate(c echo.Context) error {
	page, pageSize, limit, offset := parsePagination(c)
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("ListMaintenanceTemplate: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	queryParams := models.MaintenanceTemplateParams{}
	queryParams.Limit = limit
	queryParams.Offset = offset

	rows, total, err := dbrepo.GetMaintenanceTemplate(pool, queryParams, true)
	if err != nil {
		c.Logger().Errorf("ListMaintenanceTemplate db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query templates: %v", err)})
	}

	var out []models.MaintenanceTemplate
	for _, r := range rows {
		out = append(out, models.MaintenanceTemplate{
			UID:         r.UID,
			Name:        r.Name,
			Description: r.Description,
			TypeUID:     *r.TypeUID,
			Frequency:   r.Frequency,
		})
	}
	if out == nil {
		out = []models.MaintenanceTemplate{}
	}
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data: out,
		Pagination: models.PaginationInfo{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	})
}

// GetMaintenanceTemplate godoc
// @Summary Get Maintenance Template
// @Description Get maintenance template by UID
// @Tags MaintenanceTemplates
// @Produce json
// @Param uid path string true "Template UID"
// @Success 200 {object} models.MaintenanceTemplate
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance/templates/{uid} [get]
func GetMaintenanceTemplate(c echo.Context) error {
	// 1. Read UID from path
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}

	// 2. Get DB instance
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "db not found in context"})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("GetMaintenanceTemplate: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal error: db object invalid"})
	}

	// 3. Prepare query params (single-record mode)
	qp := models.MaintenanceTemplateParams{
		UID:    uid,
		Limit:  1,
		Offset: 0,
	}

	// 4. Call repository
	rows, _, err := dbrepo.GetMaintenanceTemplate(pool, qp, false)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("GetMaintenanceTemplate db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("db query template: %v", err),
		})
	}

	if len(rows) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
	}

	// 5. Since it's a single record request, return rows[0]
	r := rows[0]

	response := models.MaintenanceTemplate{
		UID:         r.UID,
		Name:        r.Name,
		Description: r.Description,
		TypeUID: func() string {
			if r.TypeUID != nil {
				return *r.TypeUID
			}
			return ""
		}(),
		Frequency: r.Frequency,
	}

	return c.JSON(http.StatusOK, response)
}

// CreateMaintenanceTemplate godoc
// @Summary Create Maintenance Template
// @Description Create a new maintenance template
// @Tags MaintenanceTemplates
// @Accept json
// @Produce json
// @Param template body models.MaintenanceTemplateCreateRequest true "Template payload"
// @Success 201 {object} map[string]string "{\"uid\":\"550e8400-e29b-41d4-a716-446655440000\"}"
// @Failure 400 {object} map[string]string "{\"error\":\"Invalid request\"}"
// @Failure 500 {object} map[string]string "{\"error\":\"Internal server error\"}"
// @Router /maintenance/templates [post]
func CreateMaintenanceTemplate(c echo.Context) error {
	var req models.MaintenanceTemplateCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - template payload received", "template": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("CreateMaintenanceTemplate: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	frequencyRow, err := typeRepo.GetTypeByUID(pool, req.FrequencyUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid frequency UID",
		})
	}

	queryParams := models.MaintenanceType{}
	queryParams.UID = req.TypeUID
	typeRows, err := dbrepo.GetMaintenanceType(pool, queryParams)
	typeRow := typeRows[0]
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	newUID, err := dbrepo.CreateMaintenanceTemplate(pool, req.Name, req.Description, typeRow.ID, frequencyRow.ID)
	if err != nil {
		c.Logger().Errorf("CreateMaintenanceTemplate db insert error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db insert template: %v", err)})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

// UpdateMaintenanceTemplate godoc
// @Summary Update Maintenance Template
// @Description Update maintenance template by UID
// @Tags MaintenanceTemplates
// @Accept json
// @Produce json
// @Param uid path string true "Template UID"
// @Param template body models.MaintenanceTemplateUpdateRequest true "Template update payload"
// @Success 200 {object} map[string]string "{\"uid\":\"550e8400-e29b-41d4-a716-446655440000\"}"
// @Failure 400 {object} map[string]string "{\"error\":\"Invalid request\"}"
// @Failure 500 {object} map[string]string "{\"error\":\"Internal server error\"}"
// @Router /maintenance/templates/{uid} [patch]
func UpdateMaintenanceTemplate(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var req models.MaintenanceTemplateUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would update template", "uid": uid, "payload": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("UpdateMaintenanceTemplate: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	queryParams := models.MaintenanceTemplateParams{UID: uid}

	rows, _, err := dbrepo.GetMaintenanceTemplate(pool, queryParams, false)
	if err != nil || len(rows) == 0 {
		if err == pgx.ErrNoRows || len(rows) == 0 {
			return c.JSON(http.StatusNotFound, echo.Map{"error": fmt.Sprintf("UID: %v not found", uid)})
		}
		c.Logger().Errorf("DeleteMaintenanceTemplate lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query template: %v", err)})
	}

	patch := dbrepo.MaintenanceTemplateUpdate{}

	// ---- TYPE_UID ----
	if req.TypeUID != nil && *req.TypeUID != "" {
		// validate maintenance type exists
		qp := models.MaintenanceType{UID: *req.TypeUID}
		typeRows, err := dbrepo.GetMaintenanceType(pool, qp)
		if err != nil {
			return fmt.Errorf("invalid type_uid: %v", err)
		}
		if len(typeRows) == 0 {
			return fmt.Errorf("type_uid not found")
		}
		typeRow := typeRows[0]

		patch.TypeID = &typeRow.ID

	}

	if req.FrequencyUID != nil && *req.FrequencyUID != "" {
		// validate maintenance type exists
		frequencyRow, err := typeRepo.GetTypeByUID(pool, *req.FrequencyUID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid frequency UID"})
		}

		patch.FrequencyID = &frequencyRow.ID

	}

	if req.Name != nil && *req.Name != "" {
		patch.Name = req.Name
	}

	if req.Description != nil && *req.Description != "" {
		patch.Description = req.Description
	}

	fmt.Println(patch)

	if err := dbrepo.UpdateMaintenanceTemplate(pool, uid, patch); err != nil {
		c.Logger().Errorf("UpdateMaintenanceTemplate db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db update template: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

// DeleteMaintenanceTemplate godoc
// @Summary Delete Maintenance Template
// @Description Delete a maintenance template by UID
// @Tags MaintenanceTemplates
// @Produce json
// @Param uid path string true "Template UID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance/templates/{uid} [delete]
func DeleteMaintenanceTemplate(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would not delete template", "uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteMaintenanceTemplate: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	queryParams := models.MaintenanceTemplateParams{UID: uid}

	rows, _, err := dbrepo.GetMaintenanceTemplate(pool, queryParams, false)
	if err != nil || len(rows) == 0 {
		if err == pgx.ErrNoRows || len(rows) == 0 {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("DeleteMaintenanceTemplate lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query template: %v", err)})
	}
	if err := dbrepo.DeleteMaintenanceTemplate(pool, uid); err != nil {
		c.Logger().Errorf("DeleteMaintenanceTemplate db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete template: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Template deleted successfully"})
}

func parsePagination(c echo.Context) (int, int, int, int) {
	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("page_size")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	limit := pageSize

	return page, pageSize, limit, offset
}

func GetMaintenanceTemplateByUID(pool *pgxpool.Pool, uid string) (*dbrepo.MaintenanceTemplateRow, error) {
	templateParams := models.MaintenanceTemplateParams{UID: uid}
	templateRows, _, err := dbrepo.GetMaintenanceTemplate(pool, templateParams, false)
	if err != nil || len(templateRows) == 0 {
		return nil, fmt.Errorf("template UID does not exist")
	}
	return &templateRows[0], nil
}
