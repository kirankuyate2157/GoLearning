package routes

import (
	controllers "app/controllers/settings"

	"github.com/labstack/echo/v4"
)

func SettingRoutes(e *echo.Echo) {

	g := e.Group("/api/settings")

	g.GET("/alerts/types", controllers.ListAlertsTypes)

	g.GET("/alerts", controllers.ListAlerts)
	g.GET("/alerts/:uid", controllers.ListAlerts)
	g.POST("/alerts", controllers.CreateAlerts)
	g.PATCH("/alerts/:uid", controllers.UpdateAlert)
	g.DELETE("/alerts/:uid", controllers.DeleteAlert)

	g.GET("/taxes", controllers.ListTaxes)
	g.POST("/taxes", controllers.CreateTaxes)
	g.PATCH("/taxes/:uid", controllers.UpdateTaxes)
	g.DELETE("/taxes/:uid", controllers.DeleteTaxes)

}
