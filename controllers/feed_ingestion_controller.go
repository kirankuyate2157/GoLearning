package controllers

import (
	dbrepo "app/db/repo"
	"app/ingestion"
	"os"
	"strings"
	"app/models"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// ListActiveIngestions godoc
// @Summary List Active Ingestion Processes
// @Description Get list of all active ingestion processes
// @Tags FeedIngetions
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /feeds/ingestions [get]
func ListActiveIngestions(c echo.Context) error {
	activeProcesses := ingestion.Manager.ListActiveProcesses()
	return c.JSON(http.StatusOK, echo.Map{
		"active_ingestions": activeProcesses,
		"total_count":       len(activeProcesses),
	})
}

// StartIngestion godoc
// @Summary Start Feed Ingestion Process
// @Description Start an ingestion subprocess for a specific feed
// @Tags FeedIngetions
// @Param feed_uid path string true "Feed UID"
// @Produce json
// @Success 200 {object} map[string]string "Ingestion started successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Feed not found"
// @Failure 409 {object} map[string]string "Ingestion already running"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /feeds/{feed_uid}/ingestions [post]
func StartIngestion(c echo.Context) error {
	feedUID := c.Param("uid")
	var isRunning bool = true

	if feedUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "feed_uid is required"})
	}


	// Check if already running
	if ingestion.Manager.IsProcessActive(feedUID) {
		return c.JSON(http.StatusConflict, echo.Map{"error": "ingestion process is already running for this feed"})
	}

	// Get database pool
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("StartIngestion: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	// Get feed data from database
	feedData, err := dbrepo.GetFeedByUID(pool, feedUID)
	if err != nil {
		c.Logger().Errorf("StartIngestion: failed to get feed data: %v", err)
		return c.JSON(http.StatusNotFound, echo.Map{"error": fmt.Sprintf("failed to get feed data: %v", err)})
	}

	if feedData.IsRunning{
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Ingestion process is already running"})
	}

	// Build CLI arguments for the ingestion script
	feedArgs, err := buildIngestionArgs(feedData)

	if err != nil {
		c.Logger().Errorf("Error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintln(err)})

	}

	// Start the ingestion process
	err = ingestion.Manager.StartProcess(
		feedData.UID,
		feedData.ConnectionType,
		feedArgs,
	)
	// return c.JSON(http.StatusNotImplemented, echo.Map{"error": "IMPLEMENTING OR TESTING"})

	if err != nil {
		c.Logger().Errorf("StartIngestion: failed to start process: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("failed to start ingestion: %v", err)})
	}

	var req models.FeedUpdateRequest
	req.IsRunning = &isRunning
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if err := dbrepo.UpdateFeed(pool, feedUID, req); err != nil {
		c.Logger().Errorf("UpdateFeed db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": fmt.Sprintf("db update feed: %v", err),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":         "ingestion started successfully",
		"feed_uid":        feedData.UID,
		"connection_type": feedData.ConnectionType,
	})
}





// StopIngestion godoc
// @Summary Stop Feed Ingestion Process
// @Description Stop an active ingestion subprocess for a specific feed
// @Tags FeedIngetions
// @Param uid path string true "Feed UID"
// @Produce json
// @Success 200 {object} map[string]string "Ingestion stopped successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "No active ingestion found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /feeds/{uid}/ingestions [delete]
func StopIngestion(c echo.Context) error {
	feedUID := strings.Trim(c.Param("uid"), "/")
	if feedUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "feed_uid is required"})
	}

	var (
		remark     *string
		errMessage string
		httpStatus int
	)

	// Stop the ingestion process
	if err := ingestion.Manager.StopProcess(feedUID); err != nil {
		msg := fmt.Sprintf("failed to stop ingestion: %v", err)
		remark = &msg
		errMessage = msg
		httpStatus = http.StatusNotFound
	} else {
		httpStatus = http.StatusOK
	}

	// Get DB pool
	dbi := c.Get("db")
	if dbi == nil {
		msg := "database not found in context"
		remark = &msg
		errMessage = msg
		httpStatus = http.StatusInternalServerError
	} 
	
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		msg := "invalid db object in context"
		remark = &msg
		errMessage = msg
		httpStatus = http.StatusInternalServerError
	}
	// Update feed is_running = false
	isRunning := false
	updateReq := models.FeedUpdateRequest{
		IsRunning: &isRunning,
	}
	if err := dbrepo.UpdateFeed(pool, feedUID, updateReq); err != nil {
		msg := fmt.Sprintf("db update feed failed: %v", err)
		remark = &msg
		errMessage = msg
		httpStatus = http.StatusBadRequest
	}

	// Fetch feed info for history
	feedData, err := dbrepo.GetFeedByUID(pool, feedUID)
	if err != nil {
		msg := fmt.Sprintf("failed to get feed data: %v", err)
		remark = &msg
		errMessage = msg
		httpStatus = http.StatusNotFound
	}


	// Create feed history
	if  err := dbrepo.CreateFeedHistoryEntry(pool, feedData.ID, &isRunning, remark); err != nil {
		c.Logger().Errorf("failed to save feed history: %v", err)
		if httpStatus == http.StatusOK {
			httpStatus = http.StatusInternalServerError
			errMessage = fmt.Sprintf("failed to save feed history: %v", err)
		}
	}

	// Return response
	if httpStatus != http.StatusOK && errMessage != "" {
		return c.JSON(httpStatus, echo.Map{"error": errMessage})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message":  "ingestion stopped successfully",
		"feed_uid": feedUID,
	})
}


// buildIngestionArgs constructs CLI arguments based on feed data and connection type
func buildIngestionArgs(feedData dbrepo.FeedRow) ([]string, error) {
	args := []string{}

	// Get TimescaleDB connection info from environment
	dbHost := os.Getenv("TIMESCALEDB_HOST")
	dbPort := os.Getenv("TIMESCALEDB_PORT")
	dbUser := os.Getenv("TIMESCALEDB_USER")
	dbPass := os.Getenv("TIMESCALEDB_PASSWORD")
	dbName := os.Getenv("TIMESCALEDB_DB")

	// Common arguments for all ingestion types (TimescaleDB connection)
	args = append(args,
		"--db-host", dbHost,
		"--db-port", dbPort,
		"--db-user", dbUser,
		"--db-password", dbPass,
		"--db-name", dbName,
		"--feed-uid", feedData.UID,
	)

	var fieldArgs []string
	for i := 0; i < len(feedData.JSONData); i++ {
		v := fmt.Sprintf("%v", feedData.JSONData[i]["value"])
		k := fmt.Sprintf("--%v", feedData.JSONData[i]["field-name"])

		if v == "" {
			continue
		}
		// fmt.Printf("%T => %v \n", v, v)
		fieldArgs = append(fieldArgs, k)
		fieldArgs = append(fieldArgs, v)
	}

	args = append(args, fieldArgs...)

	return args, nil
}
