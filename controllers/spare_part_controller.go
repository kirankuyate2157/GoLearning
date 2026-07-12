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

// GetSpareParts godoc
// @Summary Get Spare Parts
// @Description Get list of spare parts (JWT protected)
// @Tags SpareParts
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.SparePart
// @Router /spare-parts [get]
func GetSpareParts(c echo.Context) error {
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.SparePart{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("GetSpareParts: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	rows, err := dbrepo.ListSpareParts(pool)
	if err != nil {
		c.Logger().Errorf("GetSpareParts db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query spare_parts: %v", err)})
	}

	var out []models.SparePart
	for _, r := range rows {
		var sp models.SparePart
		if r.UID.Valid {
			sp.UID = r.UID.String
		}
		sp.Name = r.Name
		sp.Code = r.Code
		if r.Description.Valid {
			sp.Description = r.Description.String
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			sp.CreatedAt = &t
		}
		if r.TypeUID.Valid {
			sp.TypeUID = r.TypeUID.String
		}

		// Fetch files for this spare part
		entityType := "spare_part"
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
		sp.Files = files

		out = append(out, sp)
	}
	if out == nil {
		out = []models.SparePart{}
	}
	return c.JSON(http.StatusOK, out)
}

// GetSparePart godoc
// @Summary Get Spare Part
// @Description Get spare part by UID (JWT protected)
// @Tags SpareParts
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Spare Part UID"
// @Success 200 {object} models.SparePart
// @Failure 404 {object} map[string]string
// @Router /spare-parts/{uid} [get]
func GetSparePart(c echo.Context) error {
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
		c.Logger().Error("GetSparePart: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	r, err := dbrepo.GetSparePartByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("GetSparePart db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query spare_part: %v", err)})
	}

	// Fetch files for this spare part
	entityType := "spare_part"
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

	var sp models.SparePart
	if r.UID.Valid {
		sp.UID = r.UID.String
	}
	sp.Name = r.Name
	sp.Code = r.Code
	if r.Description.Valid {
		sp.Description = r.Description.String
	}
	if r.CreatedAt.Valid {
		t := r.CreatedAt.Time
		sp.CreatedAt = &t
	}
	if r.TypeUID.Valid {
		sp.TypeUID = r.TypeUID.String
	}
	sp.Files = files
	return c.JSON(http.StatusOK, sp)
}

// CreateSparePart godoc
// @Summary Create Spare Part
// @Description Create a new spare part (JWT protected)
// @Tags SpareParts
// @Accept json
// @Produce json
// @Param spare_part body models.SparePartCreateRequest true "Spare part payload"
// @Security BearerAuth
// @Router /spare-parts [post]
func CreateSparePart(c echo.Context) error {
	var req models.SparePartCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Validate required fields
	if strings.TrimSpace(req.Name) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "name is required"})
	}
	if strings.TrimSpace(req.Code) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "code is required"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - spare_part payload received", "spare_part": req})
	}
	pool := dbi.(*pgxpool.Pool)

	// validate type
	t, err := dbrepo.GetTypeByUID(pool, req.TypeUID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "invalid type_id uid"})
	}

	newUID, err := dbrepo.CreateSparePart(pool, req, &t.ID)
	if err != nil {
		c.Logger().Errorf("CreateSparePart db insert error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db insert spare_part: %v", err)})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

// UpdateSparePart godoc
// @Summary Update Spare Part
// @Description Update an existing spare part by UID (JWT protected)
// @Tags SpareParts
// @Accept json
// @Produce json
// @Param uid path string true "Spare Part UID"
// @Param spare_part body models.SparePartUpdateRequest true "Spare part update payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /spare-parts/{uid} [patch]
func UpdateSparePart(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	var req models.SparePartUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	hasUpdate := req.Name != nil || req.Code != nil || req.Description != nil
	if !hasUpdate {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "no fields to update"})
	}

	// Validate non-empty if provided
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "name cannot be empty"})
	}
	if req.Code != nil && strings.TrimSpace(*req.Code) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "code cannot be empty"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would update spare_part", "uid": uid, "payload": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("UpdateSparePart: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	var type_ *int
	if req.TypeUID != nil {
		t, err := dbrepo.GetTypeByUID(pool, *req.TypeUID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid type_id uid"})
		}
		type_ = &t.ID
	}

	if err := dbrepo.UpdateSparePart(pool, uid, req, type_); err != nil {
		c.Logger().Errorf("UpdateSparePart db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db update spare_part: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

// DeleteSparePart godoc
// @Summary Delete Spare Part
// @Description Delete a spare part by UID (JWT protected)
// @Tags SpareParts
// @Produce json
// @Param uid path string true "Spare Part UID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /spare-parts/{uid} [delete]
func DeleteSparePart(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would delete spare_part", "uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteSparePart: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	if _, err := dbrepo.GetSparePartByUID(pool, uid); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("DeleteSparePart lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query spare_part: %v", err)})
	}

	if err := dbrepo.DeleteSparePart(pool, uid); err != nil {
		c.Logger().Errorf("DeleteSparePart db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete spare_part: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}
