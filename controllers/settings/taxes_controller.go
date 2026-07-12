package controllers

import (
	dbrepo "app/db/repo/settings"
	models "app/models/settings"
	"fmt"
	"net/http"

	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// ListTaxes godoc
// @Summary Get Taxes List
// @Description Get Taxes List
// @Tags Taxes
// @Produce json
// @Success 200 {object} models.TaxesGet
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /settings/taxes [get]
func ListTaxes(c echo.Context) error {
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "database not found in context",
		})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("ListTaxes: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "invalid db object in context",
		})
	}

	rows, err := dbrepo.ListTaxes(pool)
	if err != nil {
		c.Logger().Errorf("ListTaxes db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("db query taxes: %v", err),
		})
	}

	var output []models.TaxesGet

	for _, r := range rows {
		var ft models.TaxesGet

		if r.UID.Valid {
			ft.UID = r.UID.String
		}

		if r.Type.Valid {
			ft.TYPE = r.Type.String
		}

		ft.PERCENTAGE = r.Percent

		output = append(output, ft)
	}

	// Always return array, never null
	if output == nil {
		output = []models.TaxesGet{}
	}

	return c.JSON(http.StatusOK, output)
}

// CreateTaxes godoc
// @Summary Create Tax
// @Description Create a new tax
// @Tags Taxes
// @Accept json
// @Produce json
// @Param taxes body models.TaxesPost true "Tax payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /settings/taxes [post]
func CreateTaxes(c echo.Context) error {
	var req models.TaxesPost
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"note":  "no db configured - tax payload received",
			"taxes": req,
		})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "invalid db object in context",
		})
	}

	repoReq := dbrepo.TaxesCreate{
		Type:       req.Type,
		Percentage: req.Percentage,
	}

	newUID, err := dbrepo.CreateTaxes(pool, repoReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": fmt.Sprintf("db insert taxes: %v", err),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

// UpdateTaxes godoc
// @Summary Update Tax
// @Description Update tax by UID
// @Tags Taxes
// @Accept json
// @Produce json
// @Param uid path string true "Tax UID"
// @Param taxes body models.TaxesPatch true "Tax update payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /settings/taxes/{uid} [patch]
func UpdateTaxes(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}
	var req models.TaxesPatch
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would update tax", "uid": uid, "payload": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("UpdateTaxes: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	// Check tax exists
	taxList, err := dbrepo.ListTaxes(pool)
	if err != nil {
		c.Logger().Errorf("ListTaxes db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db query taxes error"})
	}

	if len(taxList) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Tax not found"})
	}

	taxID := taxList[0].ID

	data := dbrepo.TaxesUpdate{
		Type:       req.Type,
		Percentage: req.Percentage,
	}

	if err := dbrepo.UpdateTaxes(pool, taxID, data); err != nil {
		c.Logger().Errorf("UpdateTaxes db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db update tax: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})

}

// DeleteTaxes godoc
// @Summary Delete Tax
// @Description Delete a tax by UID
// @Tags Taxes
// @Produce json
// @Param uid path string true "Tax UID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /settings/taxes/{uid} [delete]
func DeleteTaxes(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")

	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid tax uid"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would not delete tax", "uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteTaxes: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	// Check tax exists
	taxList, err := dbrepo.ListTaxes(pool)
	if err != nil {
		c.Logger().Errorf("ListTaxes db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db query taxes error"})
	}

	if len(taxList) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Tax not found"})
	}

	id := taxList[0].ID

	if err := dbrepo.DeleteTaxes(pool, id); err != nil {
		c.Logger().Errorf("DeleteTaxes db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete tax: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "tax deleted successfully"})
}
