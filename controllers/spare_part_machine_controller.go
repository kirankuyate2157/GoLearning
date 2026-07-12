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

// Helper assumed functions in your repo:
// - dbrepo.GetSparePartByUID(pool, uid) -> returns spare part row (with ID)
// - dbrepo.GetMachineByUID(pool, uid) -> returns machine row (with ID)
// If your machine repo function has a different name, adjust accordingly.

// GetSparePartMachineList godoc
// @Summary Get spare-part -> machine list
// @Tags SparePartMachines
// @Produce json
// @Param sparepart_uid path string true "Spare Part UID"
// @Success 200 {array} models.SparePartMachine
// @Router /spare-parts/{sparepart_uid}/machines [get]
func GetSparePartMachineList(c echo.Context) error {
	spareUID := strings.Trim(c.Param("sparepart_uid"), "/")
	if spareUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid spare part uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.SparePartMachine{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("GetSparePartMachineList: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	// verify spare part exists
	if _, err := dbrepo.GetSparePartByUID(pool, spareUID); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "spare part not found"})
		}
		c.Logger().Errorf("GetSparePartMachineList spare lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query spare_part: %v", err)})
	}

	rows, err := dbrepo.ListSparePartMachines(pool, spareUID)
	if err != nil {
		c.Logger().Errorf("ListSparePartMachines db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query spare_parts_machines: %v", err)})
	}

	var out []models.SparePartMachine
	for _, r := range rows {
		var spm models.SparePartMachine
		if r.UID.Valid {
			spm.UID = r.UID.String
		}
		if r.SparePartUID.Valid {
			spm.SparePartUID = r.SparePartUID.String
		}
		if r.MachineUID.Valid {
			spm.MachineUID = r.MachineUID.String
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			spm.CreatedAt = &t
		}
		out = append(out, spm)
	}
	if out == nil {
		out = []models.SparePartMachine{}
	}
	return c.JSON(http.StatusOK, out)
}

// CreateSparePartMachine godoc
// @Summary Create spare-part -> machine mapping
// @Tags SparePartMachines
// @Accept json
// @Produce json
// @Param sparepart_uid path string true "Spare Part UID"
// @Param payload body models.SparePartMachineCreateRequest true "Mapping payload"
// @Success 201 {object} map[string]string
// @Router /spare-parts/{sparepart_uid}/machines [post]
func CreateSparePartMachine(c echo.Context) error {
	spareUID := strings.Trim(c.Param("sparepart_uid"), "/")
	if spareUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid spare part uid"})
	}

	var req models.SparePartMachineCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if strings.TrimSpace(req.MachineUID) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "machine_uid is required"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - would create mapping", "payload": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("CreateSparePartMachine: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	// resolve spare part id
	spareRow, err := dbrepo.GetSparePartByUID(pool, spareUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "spare part not found"})
		}
		c.Logger().Errorf("CreateSparePartMachine spare lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query spare_part: %v", err)})
	}

	// resolve machine id - assume repo has GetMachineByUID
	machineRow, err := dbrepo.GetMachineByUID(pool, req.MachineUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "machine not found"})
		}
		c.Logger().Errorf("CreateSparePartMachine machine lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine: %v", err)})
	}

	newUID, err := dbrepo.CreateSparePartMachine(pool, spareRow.ID, machineRow.ID)
	if err != nil {
		c.Logger().Errorf("CreateSparePartMachine db insert error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db insert spare_parts_machines: %v", err)})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

// DeleteSparePartMachine godoc
// @Summary Delete spare-part -> machine mapping
// @Tags SparePartMachines
// @Produce json
// @Param sparepart_uid path string true "Spare Part UID"
// @Param sparepartmachine_uid path string true "Machine UID"
// @Success 204 {object} nil
// @Router /spare-parts/{sparepart_uid}/machines/{sparepartmachine_uid} [delete]
func DeleteSparePartMachine(c echo.Context) error {
	sparePartUID := c.Param("sparepart_uid")
	machineUID := c.Param("sparepartmachine_uid")

	if sparePartUID == "" || machineUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing parameters"})
	}

	pool := c.Get("db").(*pgxpool.Pool)

	sparePart, err := dbrepo.GetSparePartByUID(pool, sparePartUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	machine, err := dbrepo.GetMachineByUID(pool, machineUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	affected, err := dbrepo.DeleteSparePartMachine(pool, sparePart.ID, machine.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if affected == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "mapping not found"})
	}

	return c.NoContent(http.StatusNoContent)
}
