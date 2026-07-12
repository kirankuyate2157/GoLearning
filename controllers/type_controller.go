package controllers

import (
	dbrepo "app/db/repo"
	"app/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// ListTypes godoc
// @Summary Get Types
// @Description Get list of types with optional category filter
// @Tags Types
// @Produce json
// @Param category query string false "Filter by category"
// @Success 200 {array} models.Type
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /types [get]
func ListTypes(c echo.Context) error {
	category := strings.TrimSpace(c.QueryParam("category"))

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("ListTypes: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	rows, err := dbrepo.ListTypes(pool, category)
	if err != nil {
		c.Logger().Errorf("ListTypes db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("db query types: %v", err),
		})
	}

	var out []models.Type
	for _, r := range rows {
		var t models.Type
		if r.UID.Valid {
			t.UID = r.UID.String
		}
		if r.Name.Valid {
			t.Name = r.Name.String
		}
		if r.Description.Valid {
			t.Description = r.Description.String
		}
		if r.Category.Valid {
			t.Category = r.Category.String
		}
		out = append(out, t)
	}

	if out == nil {
		out = []models.Type{}
	}

	return c.JSON(http.StatusOK, out)
}

// ListTypesCategory godoc
// @Summary Get Types Categories
// @Description Get list of types categories
// @Tags Types
// @Produce json
// @Success 200 {array} string
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /types/category [get]
func ListTypeCategory(c echo.Context) error {

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("ListTypeCategory: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	categories, err := dbrepo.ListTypeCategory(pool)
	if err != nil {
		c.Logger().Errorf("ListTypeCategory db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("db query types categories: %v", err),
		})
	}

	return c.JSON(http.StatusOK, categories)
}

// GetMachine godoc
// @Summary Get Types
// @Description Get machine by UID (JWT protected)
// @Tags Types
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Machine UID"
// @Success 200 {object} models.Machine
// @Failure 404 {object} map[string]string
// @Router /types/{uid} [get]
func GetType(c echo.Context) error {
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
		c.Logger().Error("GetMachine: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	r, err := dbrepo.GetTypeByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("GetMachine db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query types: %v", err)})
	}
	var t models.Type
	if r.UID.Valid {
		t.UID = r.UID.String
	}
	t.Name = r.Name.String
	if r.Description.Valid {
		t.Description = r.Description.String
	}
	if r.Category.Valid {
		t.Category = r.UID.String
	}
	// if r.CreatedAt.Valid {
	// 	t := r.CreatedAt.Time
	// 	t.CreatedAt = &t
	// }
	// return c.NoContent(http.StatusNoContent)
	return c.JSON(http.StatusOK, t)

}

// CreateType godoc
// @Summary Create Type
// @Description Create a new type (JWT protected)
// @Tags Types
// @Accept json
// @Produce json
// @Param type body models.TypeCreateRequest true "Type payload"
// @Security BearerAuth
// @Router /types [post]
func CreateType(c echo.Context) error {
	var req models.TypeCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - machine payload received", "machine": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)

	if !ok {
		c.Logger().Error("CreateMachine: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	newUID, err := dbrepo.CreateType(pool, req)
	if err != nil {
		c.Logger().Errorf("CreateMachine db insert error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db insert machine: %v", err)})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

// UpdateType godoc
// @Summary Update Type
// @Description Update type (JWT protected)
// @Tags Types
// @Accept json
// @Produce json
// @Param uid path string true "type UID"
// @Param type body models.TypeUpdateRequest true "Type update payload"
// @Security BearerAuth
// @Router /types/{uid} [patch]
func UpdateType(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var req models.TypeUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would update machine", "uid": uid, "payload": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("UpdateMachine: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	if err := dbrepo.UpdateType(pool, uid, req); err != nil {
		c.Logger().Errorf("UpdateMachine db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db update machine: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

// DeleteType godoc
// @Summary Delete type
// @Description Delete a type by UID (JWT protected)
// @Tags Types
// @Produce json
// @Param uid path string true "type UID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /types/{uid} [delete]
func DeleteType(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would not delete type", "uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteType: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	if _, err := dbrepo.GetTypeByUID(pool, uid); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("DeleteType lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine: %v", err)})
	}

	if err := dbrepo.DeleteType(pool, uid); err != nil {
		c.Logger().Errorf("DeleteType db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete machine: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Type deleted successfully",
	})
}
