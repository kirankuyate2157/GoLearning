package controllers

import (
	"net/http"
	"strconv"

	"app/db/repo"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type auditController struct{}

func NewAuditController() *auditController {
	return &auditController{}
}

func (a *auditController) ListAuditEventsByUser(c echo.Context) error {
	userUID := c.Param("user_uid")
	limitStr := c.QueryParam("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	ctx := c.Request().Context()
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "missing db pool in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db pool in context"})
	}

	ar := repo.NewAuditRepo(pool)
	events, err := ar.ListAuditEventsByUser(ctx, userUID, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, events)
}

func (a *auditController) ListAuditEventsByResource(c echo.Context) error {
	resourceType := c.QueryParam("type")
	resourceUID := c.QueryParam("uid")
	limitStr := c.QueryParam("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if resourceType == "" || resourceUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "type and uid query parameters required"})
	}

	ctx := c.Request().Context()
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "missing db pool in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db pool in context"})
	}

	ar := repo.NewAuditRepo(pool)
	events, err := ar.ListAuditEventsByResource(ctx, resourceType, resourceUID, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, events)
}
