package routes

import (
	"app/controllers"

	"github.com/labstack/echo/v4"
)

func HealthRoutes(e *echo.Echo) {
	g := e.Group("/api/health")
	g.GET("", controllers.HealthCheck)
	g.GET("/ready", controllers.ReadinessCheck)
	g.GET("/version", controllers.VersionInfo)
}
