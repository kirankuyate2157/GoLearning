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

// ListFeeds godoc
// @Summary Get Feeds
// @Description Get list of feeds
// @Tags Feeds
// @Produce json
// @Success 200 {array} models.Feed
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /feeds [get]
func ListFeeds(c echo.Context) error {

	dbi := c.Get("db")

	if dbi == nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("ListFeeds: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	rows, err := dbrepo.ListFeeds(pool)
	if err != nil {

		c.Logger().Errorf("ListFeeds db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query feeds: %v", err)})
	}

	out := make([]models.Feed, 0, len(rows))

	for _, r := range rows {

		f := models.Feed{
			UID:               r.UID,
			Name:              r.Name,
			Code:              r.Code,
			Description:       r.Description,
			LowValue:          r.LowValue,
			HighValue:         r.HighValue,
			DisplayUnitUID:    r.DisplayUnit,
			FeedTypeUID:       r.TypeUID,
			MarginOfError:     r.MarginOfError,
			ConnectionTypeUID: r.ConnectionTypeUID,
			IngestionType:     r.IngestionType,
			PollingRate:       r.PollingRate,
			IsRunning:         r.IsRunning,
			JSONData:          r.JSONData,
		}

		out = append(out, f)
	}

	// No need for a nil check — JSON will encode nil slice as []
	fmt.Println("ListFeeds Controller: response sent successfully")
	return c.JSON(http.StatusOK, out)
}

// GetFeed godoc
// @Summary Get Feed
// @Description Get feed by UID
// @Tags Feeds
// @Produce json
// @Param uid path string true "Feed UID"
// @Success 200 {object} models.Feed
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feeds/{uid} [get]
func GetFeed(c echo.Context) error {
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
		c.Logger().Error("GetFeed: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	r, err := dbrepo.GetFeedByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("GetFeed db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query feeds: %v", err)})
	}
	var f models.Feed
	f.UID = r.UID
	f.Name = r.Name
	f.Code = r.Code
	f.Description = r.Description
	f.LowValue = r.LowValue
	f.HighValue = r.HighValue
	f.DisplayUnitUID = r.DisplayUnit
	f.FeedTypeUID = r.TypeUID
	f.MarginOfError = r.MarginOfError
	f.ConnectionTypeUID = r.ConnectionTypeUID
	f.IngestionType = r.IngestionType
	f.PollingRate = r.PollingRate
	f.IsRunning = r.IsRunning
	f.JSONData = r.JSONData
	return c.JSON(http.StatusOK, f)
}

// CreateFeed godoc
// @Summary Create Feed
// @Description Create a new feed
// @Tags Feeds
// @Accept json
// @Produce json
// @Param feed body models.FeedCreateRequest true "Feed payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feeds [post]
func CreateFeed(c echo.Context) error {

	var req models.FeedCreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "no db configured"})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object"})
	}

	// --------------------------------------------------
	// Convert display_unit_uid → display_unit_id (INT)
	// --------------------------------------------------
	du, err := dbrepo.GetTypeByUID(pool, req.DisplayUnitUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid display_unit"})
	}
	displayUnitID := du.ID

	ft, err := dbrepo.GetTypeByUID(pool, req.FeedTypeUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid type_id"})
	}
	feedTypeID := ft.ID

	// --------------------------------------------------
	// Convert connection_type_uid → connection_type_id (INT)
	// --------------------------------------------------
	ct, err := dbrepo.GetFeedTypeByUID(pool, req.ConnectionTypeUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid connection_type_uid"})
	}
	connectionTypeID := ct.ID

	// --------------------------------------------------
	// Call Repo Method
	// --------------------------------------------------
	newUID, err := dbrepo.CreateFeed(
		pool,
		req.Name,
		req.Code,
		req.Description,
		req.LowValue,
		req.HighValue,
		displayUnitID,
		feedTypeID,
		req.MarginOfError,
		connectionTypeID,
		req.IngestionType,
		req.PollingRate,
		req.JSONData,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("db insert failed: %v", err),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

type FeedMapRequest struct {
	FeedUID string `json:"feed_uid" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// MapFeedToMachine godoc
// @Summary Map Feed to Machine
// @Description Map a feed to a machine by their UIDs
// @Tags Feeds
// @Accept json
// @Produce json
// @Param uid path string true "Machine UID"
// @Param data body FeedMapRequest true "Feed request payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feeds/machine/{uid} [post]
func MapFeedToMachine(c echo.Context) error {
	var req FeedMapRequest

	machineUID := strings.Trim(c.Param("uid"), "/")
	if machineUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid machine_uid"})
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid body"})
	}
	if req.FeedUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "feed_uid required"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - feed map payload received", "feed_map_machine": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("MapFeedToMachine: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	feed, err := dbrepo.GetFeedByUID(pool, req.FeedUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "feed uid not found"})
		}
		c.Logger().Errorf("MapFeedToMachine lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query feed: %v", err)})
	}
	feedID := feed.ID
	machine, err := dbrepo.GetMachineByUID(pool, machineUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "machine uid not found"})
		}
		c.Logger().Errorf("GetMachine db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine: %v", err)})
	}
	machineID := machine.ID

	// Convert int IDs to string if required by repo
	if err := dbrepo.CreateMachineFeedMapping(pool,
		machineID,
		feedID); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message":     "mapping created",
		"machine_uid": machineUID,
		"feed_uid":    req.FeedUID,
	})
}

// UpdateFeed godoc
// @Summary Update Feed
// @Description Update feed by UID
// @Tags Feeds
// @Accept json
// @Produce json
// @Param uid path string true "Feed UID"
// @Param feed body models.FeedUpdateRequest true "Feed update payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feeds/{uid} [patch]
func UpdateFeed(c echo.Context) error {

	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	var req models.FeedUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{
			"note":    "no db configured - would update feed",
			"uid":     uid,
			"payload": req,
		})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("UpdateFeed: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "invalid db object in context",
		})
	}

	// ----------------------------------------------------------
	//  Convert connection_type_uid → connection_type_id (INT)
	// ----------------------------------------------------------
	if req.ConnectionTypeUID != nil && *req.ConnectionTypeUID != "" {

		ft, err := dbrepo.GetFeedTypeByUID(pool, *req.ConnectionTypeUID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "invalid connection_type_uid",
			})
		}

		req.ConnectionTypeID = &ft.ID
	}

	// ----------------------------------------------------------
	//  Convert display_unit_uid → display_unit_id (INT)
	// ----------------------------------------------------------
	if req.DisplayUnitUID != nil && *req.DisplayUnitUID != "" {

		du, err := dbrepo.GetTypeByUID(pool, *req.DisplayUnitUID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "invalid display_unit_uid",
			})
		}

		req.DisplayUnitID = &du.ID
	}

	if req.FeedTypeUID != nil && *req.FeedTypeUID != "" {
		ft, err := dbrepo.GetTypeByUID(pool, *req.FeedTypeUID)
		if err != nil {
			return c.JSON(http.StatusBadGateway, echo.Map{
				"error": "invalid type_id",
			})
		}
		req.FeedTypeID = &ft.ID
	}

	// ----------------------------------------------------------
	//  Perform update
	// ----------------------------------------------------------
	if err := dbrepo.UpdateFeed(pool, uid, req); err != nil {
		c.Logger().Errorf("UpdateFeed db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": fmt.Sprintf("db update feed: %v", err),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

// DeleteFeed godoc
// @Summary Delete Feed
// @Description Delete a feed by UID
// @Tags Feeds
// @Produce json
// @Param uid path string true "Feed UID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feeds/{uid} [delete]
func DeleteFeed(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would not delete feed", "uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteFeed: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	if _, err := dbrepo.GetFeedByUID(pool, uid); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("DeleteFeed lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query feed: %v", err)})
	}
	if err := dbrepo.DeleteFeed(pool, uid); err != nil {
		c.Logger().Errorf("DeleteFeed db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete feed: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Feed deleted successfully"})
}
