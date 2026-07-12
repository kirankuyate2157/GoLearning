package controllers

import (
	dbrepo "app/db/repo"
	settings_dbrepo "app/db/repo/settings"
	"app/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// ListFeedAlerts godoc
// @Summary Get Feed Alerts
// @Description Get list of feed alerts
// @Tags FeedAlerts
// @Produce json
// @Success 200 {array} models.FeedAlertGet
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /feeds/alert [get]
func ListFeedsAlert(c echo.Context) error {
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("ListFeedAlerts: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	FeedAlertsParams := dbrepo.FeedAlertsParams{}
	rows, err := dbrepo.ListFeedAlerts(pool, FeedAlertsParams)
	if err != nil {
		c.Logger().Errorf("ListFeedAlerts db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query feed alerts: %v", err)})
	}
	var out []models.FeedAlertGet
	for _, r := range rows {
		out = append(out, models.FeedAlertGet{
			UID:            r.UID,
			AlertUID:       r.AlertUID,
			LowAlertValue:  r.LowAlertValue,
			HighAlertValue: r.HighAlertValue,
			FeedUID:        r.FeedUID,
		})
	}
	if out == nil {
		out = []models.FeedAlertGet{}
	}
	return c.JSON(http.StatusOK, out)
}

// GetFeedAlert godoc
// @Summary Get Feed Alert
// @Description Get feed alert by UID
// @Tags FeedAlerts
// @Produce json
// @Param feed_uid path string true "Feed UID"
// @Success 200 {object} models.FeedAlertGet
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feeds/{feed_uid}/alert [get]
func GetFeedsAlert(c echo.Context) error {
	feed_uid := strings.Trim(c.Param("feed_uid"), "/")
	if feed_uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("GetFeedAlert: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	FeedData, err := dbrepo.GetFeedByUID(pool, feed_uid)
	if err != nil {
		c.Logger().Errorf("CreateFeedAlert feed lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query feeds: %v", err)})
	}
	feedID := FeedData.ID

	FeedAlertsParams := dbrepo.FeedAlertsParams{
		FeedID: &feedID,
	}
	rows, err := dbrepo.ListFeedAlerts(pool, FeedAlertsParams)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("GetFeedAlert db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query feed alert: %v", err)})
	}

	if len(rows) == 0 {
		return c.JSON(http.StatusOK, []models.FeedAlertGet{})
	}
	
	resp := make([]models.FeedAlertGet, 0, len(rows))
	for _, r := range rows {
		resp = append(resp, models.FeedAlertGet{
			UID:            r.UID,
			AlertUID:       r.AlertUID,
			LowAlertValue:  r.LowAlertValue,
			HighAlertValue: r.HighAlertValue,
			FeedUID:        r.FeedUID,
		})
	}

	return c.JSON(http.StatusOK, resp)
}

// CreateFeedAlert godoc
// @Summary Create Feed Alert
// @Description Create a new feed alert
// @Tags FeedAlerts
// @Accept json
// @Produce json
// @Param feed_uid path string true "Feed UID"
// @Param feed_alert body models.FeedAlertPost true "FeedAlert payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feeds/{feed_uid}/alert [post]
func CreateFeedsAlert(c echo.Context) error {
	feed_uid := strings.Trim(c.Param("feed_uid"), "/")
	if feed_uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid feed_uid"})
	}
	var req models.FeedAlertPost
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - feed alert payload received", "feed_alert": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("CreateFeedAlert: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	alertsParams := settings_dbrepo.AlertsFilters{}
	alertsParams.UID = &req.AlertUID
	alertsList, err := settings_dbrepo.ListAlerts(pool, alertsParams)
	if err != nil {
		c.Logger().Errorf("CreateFeedAlert alert type lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query alert type: %v", err)})
	}
	if len(alertsList) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid alert uid"})
	}

	FeedData, err := dbrepo.GetFeedByUID(pool, feed_uid)
	if err != nil {
		c.Logger().Errorf("CreateFeedAlert feed lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query feeds: %v", err)})
	}
	feedID := FeedData.ID

	FeedAlertData := dbrepo.FeedAlertCreateRequest{}
	FeedAlertData.AlertID = alertsList[0].ID
	FeedAlertData.LowAlertValue = req.LowAlertValue
	FeedAlertData.HighAlertValue = req.HighAlertValue
	FeedAlertData.FeedID = feedID

	newUID, err := dbrepo.CreateFeedAlert(pool, FeedAlertData)
	if err != nil {
		c.Logger().Errorf("CreateFeedAlert db insert error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db insert feed alert: %v", err)})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

// UpdateFeedAlert godoc
// @Summary Update Feed Alert
// @Description Update feed alert by UID
// @Tags FeedAlerts
// @Accept json
// @Produce json
// @Param feed_uid path string true "Feed UID"
// @Param uid path string true "FeedAlert UID"
// @Param feed_alert body models.FeedAlertPatch true "FeedAlert update payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feeds/{feed_uid}/alert/{uid} [patch]
func UpdateFeedsAlert(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	feed_uid := strings.Trim(c.Param("feed_uid"), "/")
	if feed_uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid feed_uid"})
	}
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}
	var req models.FeedAlertPatch
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would update feed alert", "uid": uid, "payload": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("UpdateFeedAlert: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	if _, err := dbrepo.GetFeedAlertByUID(pool, uid); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("UpdateFeedsAlert lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query feed alert: %v", err)})
	}

	alertsParams := settings_dbrepo.AlertsFilters{}
	alertsParams.UID = req.AlertUID
	alertsList, err := settings_dbrepo.ListAlerts(pool, alertsParams)
	if err != nil {
		c.Logger().Errorf("CreateFeedAlert alert type lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query alert type: %v", err)})
	}
	if len(alertsList) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid alert uid"})
	}

	FeedAlertData := dbrepo.FeedAlertUpdateRequest{}
	FeedAlertData.AlertID = &alertsList[0].ID
	FeedAlertData.LowAlertValue = req.LowAlertValue
	FeedAlertData.HighAlertValue = req.HighAlertValue

	if err := dbrepo.UpdateFeedAlert(pool, uid, FeedAlertData); err != nil {
		c.Logger().Errorf("UpdateFeedAlert db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db update feed alert: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

// DeleteFeedAlert godoc
// @Summary Delete Feed Alert
// @Description Delete a feed alert by UID
// @Tags FeedAlerts
// @Produce json
// @Param uid path string true "FeedAlert UID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feeds/alert/{uid} [delete]
func DeleteFeedsAlert(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would not delete feed alert", "uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteFeedAlert: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	if _, err := dbrepo.GetFeedAlertByUID(pool, uid); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("DeleteFeedAlert lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query feed alert: %v", err)})
	}
	if err := dbrepo.DeleteFeedAlert(pool, uid); err != nil {
		c.Logger().Errorf("DeleteFeedAlert db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete feed alert: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Feed alert deleted successfully"})
}
