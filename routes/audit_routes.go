package routes

import (
	"app/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterAuditRoutes(e *echo.Echo) {
	ctl := controllers.NewAuditController()
	g := e.Group("/api/audit")

	g.GET("/users/:user_uid/events", ctl.ListAuditEventsByUser)
	g.GET("/resources/events", ctl.ListAuditEventsByResource)
}
