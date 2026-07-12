package controllers

import (
	dbrepo "app/db/repo"
	"app/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// ListFeedMachineMappings godoc
// @Summary List machine mappings for a feed
// @Description Returns mapping_uid, feed_uid, machine_uid
// @Tags Feed-Machine Mapping
// @Param feed_uid path string true "Feed UID"
// @Produce json
// @Success 200 {array} models.FeedMachineMappingResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feeds/{feed_uid}/machines [get]
func ListFeedMachine(c echo.Context) error {
	feedUID := strings.Trim(c.Param("feed_uid"), "/")
	if feedUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid feed uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not configured"})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	rows, err := dbrepo.ListFeedMachineMappings(pool, feedUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("db query mappings: %v", err),
		})
	}

	if rows == nil {
		rows = []models.FeedMachineMappingResponse{}
	}

	return c.JSON(http.StatusOK, rows)
}

// ListMachineFeed godoc
// @Summary List feed mappings for a machine
// @Description Returns mapping_uid, feed_uid, machine_uid
// @Tags Feed-Machine Mapping
// @Param machine_uid path string true "Machine UID"
// @Produce json
// @Success 200 {array} models.FeedMachineMappingResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feeds/machines/{machine_uid}/ [get]
func ListMachineFeed(c echo.Context) error {

	machineUID := strings.Trim(c.Param("machine_uid"), "/")
	if machineUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid machine uid",
		})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "database not configured",
		})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "invalid db object in context",
		})
	}

	rows, err := dbrepo.ListMachineFeedMappings(pool, machineUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("db query mappings: %v", err),
		})
	}

	if rows == nil {
		rows = []models.FeedMachineMappingResponse{}
	}

	return c.JSON(http.StatusOK, rows)
}

// AttachFeedToMachine godoc
// @Summary Attach a feed to a machine
// @Description Creates a mapping between a feed and a machine. Prevents duplicate mappings.
// @Tags Feed-Machine Mapping
// @Accept json
// @Produce json
// @Param feed_uid path string true "Feed UID"
// @Param machine_uid path string true "Machine UID"
// @Success 200 {object} map[string]interface{} "Mapping created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or mapping already exists"
// @Failure 404 {object} map[string]interface{} "Feed or Machine not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /feeds/{feed_uid}/machines/{machine_uid} [post]
func CreateFeedMachine(c echo.Context) error {
	// var req models.CreateFeedMachineRequest
	feedUID := strings.Trim(c.Param("feed_uid"), "/")
	machineUID := strings.Trim(c.Param("machine_uid"), "/")

	if feedUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid feed_uid"})
	}

	if machineUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid machine_uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"note":        "no db configured - feed payload received",
			"feed_uid":    feedUID,
			"machine_uid": machineUID,
		})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("CreateFeed: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "invalid db object in context",
		})
	}

	feed, err := dbrepo.GetFeedByUID(pool, feedUID)

	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "feed not found",
		})
	}
	machine, err := dbrepo.GetMachineByUID(pool, machineUID)

	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "machine not found",
		})
	}

	exists, err := dbrepo.FeedMachineMappingExists(pool, feed.ID, machine.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to check mapping",
		})
	}

	if exists {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "mapping already exists",
		})
	}

	err = dbrepo.CreateFeedMachineMapping(pool, feed.ID, machine.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": fmt.Sprintf("failed to attach feed: %v", err),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "Feed successfully mapped to Machine"})

}

// RemoveFeedFromMachine godoc
// @Summary Remove feed mapping from machine
// @Description Deletes a mapping between a feed and a machine if it exists.
// @Tags Feed-Machine Mapping
// @Accept json
// @Produce json
// @Param feed_uid path string true "Feed UID"
// @Param machine_uid path string true "Machine UID"
// @Success 200 {object} map[string]interface{} "Mapping deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Mapping not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /feeds/{feed_uid}/machines/{machine_uid} [delete]
func RemoveFeedFromMachine(c echo.Context) error {
	feedUID := strings.Trim(c.Param("feed_uid"), "/")
	machineUID := strings.Trim(c.Param("machine_uid"), "/")

	if feedUID == "" || machineUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db not available"})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object"})
	}

	// Fetch feed ID
	feed, err := dbrepo.GetFeedByUID(pool, feedUID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "feed not found"})
	}

	// Fetch machine ID
	machine, err := dbrepo.GetMachineByUID(pool, machineUID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "machine not found"})
	}

	// Check if mapping exists
	exists, err := dbrepo.FeedMachineMappingExists(pool, feed.ID, machine.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to check mapping"})
	}

	if !exists {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "mapping does not exist"})
	}

	// Delete mapping
	err = dbrepo.DeleteFeedMachineMapping(pool, feed.ID, machine.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to delete mapping"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "feed successfully removed from machine",
	})
}
