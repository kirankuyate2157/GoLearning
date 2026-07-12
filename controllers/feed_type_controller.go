package controllers

import (
	dbrepo "app/db/repo"
	"app/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// ListFeedTypes godoc
// @Summary Get FeedTypes
// @Description Get list of feed_types
// @Tags FeedTypes
// @Produce json
// @Success 200 {array} models.FeedType
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /feeds/type [get]
func ListFeedTypes(c echo.Context) error {
	// Get DB from context
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "database not found in context",
		})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("ListFeedTypes: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "invalid db object in context",
		})
	}

	// Fetch rows from DB
	rows, err := dbrepo.ListFeedTypes(pool)
	if err != nil {
		c.Logger().Errorf("ListFeedTypes db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("db query feed types: %v", err),
		})
	}

	// Final output
	var output []models.FeedType

	// Loop through DB rows
	for _, r := range rows {
		var ft models.FeedType

        // UID
		if r.UID.Valid {
			ft.UID = r.UID.String
		}

        // Name
		if r.Name.Valid {
			ft.Name = r.Name.String
		}

        // Description
		if r.Description.Valid {
			ft.Description = r.Description.String
		}

		// -------- JSON UNMARSHAL FOR STRUCTURE --------
		if r.Strucuture.Valid && r.Strucuture.String != "" {
			var jsonMap []map[string]interface{}

			if err := json.Unmarshal([]byte(r.Strucuture.String), &jsonMap); err != nil {
				c.Logger().Errorf("failed to unmarshal structure JSON for UID %s: %v", ft.UID, err)
				ft.Strucuture = []map[string]interface{}{} // fallback
			} else {
				ft.Strucuture = jsonMap
			}

		} else {
			ft.Strucuture = []map[string]interface{}{}
		}
		// -------------------------------------------------

		output = append(output, ft)
	}

	// Always return array, never null
	if output == nil {
		output = []models.FeedType{}
	}

	return c.JSON(http.StatusOK, output)
}







// GetFeedType godoc
// @Summary Get FeedType
// @Description Get a feed type by UID
// @Tags FeedTypes
// @Produce json
// @Param uid path string true "FeedType UID"
// @Success 200 {object} models.FeedType
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feeds/type/{uid} [get]
func GetFeedType(c echo.Context) error {
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
		c.Logger().Error("GetFeedType: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	// Fetch record from DB
	r, err := dbrepo.GetFeedTypeByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("GetFeedType db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query feed types: %v", err)})
	}

	// Convert DB row to model
	var ft models.FeedType

	if r.UID.Valid {
		ft.UID = r.UID.String
	}
	if r.Name.Valid {
		ft.Name = r.Name.String
	}
	if r.Description.Valid {
		ft.Description = r.Description.String
	}

	// -------- FIXED STRUCTURE JSON UNMARSHAL --------
	if r.Strucuture.Valid && r.Strucuture.String != "" {
		var jsonMap []map[string]interface{}
		if err := json.Unmarshal([]byte(r.Strucuture.String), &jsonMap); err != nil {
			c.Logger().Errorf("failed to unmarshal structure JSON for UID %s: %v", ft.UID, err)
			ft.Strucuture = []map[string]interface{}{} // fallback empty
		} else {
			ft.Strucuture = jsonMap
		}
	} else {
		ft.Strucuture = []map[string]interface{}{}
	}
	// -------------------------------------------------

	return c.JSON(http.StatusOK, ft)
}


// CreateFeedType godoc
// @Summary Create FeedType
// @Description Create a new feed type
// @Tags FeedTypes
// @Accept json
// @Produce json
// @Param feed_type body models.FeedTypeCreateRequest true "FeedType payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feeds/type [post]
func CreateFeedType(c echo.Context) error {
	var req models.FeedTypeCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - feed type payload received", "feed_type": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("CreateFeedType: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	newUID, err := dbrepo.CreateFeedType(pool, req)
	if err != nil {
		c.Logger().Errorf("CreateFeedType db insert error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db insert feed type: %v", err)})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

// UpdateFeedType godoc
// @Summary Update FeedType
// @Description Update an existing feed type by UID
// @Tags FeedTypes
// @Accept json
// @Produce json
// @Param uid path string true "FeedType UID"
// @Param feed_type body models.FeedTypeUpdateRequest true "FeedType update payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feeds/type/{uid} [patch]
func UpdateFeedType(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var req models.FeedTypeUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would update feed type", "uid": uid, "payload": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("UpdateFeedType: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	if err := dbrepo.UpdateFeedType(pool, uid, req); err != nil {
		c.Logger().Errorf("UpdateFeedType db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db update feed type: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

// DeleteFeedType godoc
// @Summary Delete FeedType
// @Description Delete a feed type by UID
// @Tags FeedTypes
// @Produce json
// @Param uid path string true "FeedType UID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feeds/type/{uid} [delete]
func DeleteFeedType(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would not delete feed type", "uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteFeedType: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	if _, err := dbrepo.GetFeedTypeByUID(pool, uid); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("DeleteFeedType lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query feed type: %v", err)})
	}

	if err := dbrepo.DeleteFeedType(pool, uid); err != nil {
		c.Logger().Errorf("DeleteFeedType db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete feed type: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Feed type deleted successfully"})
}
