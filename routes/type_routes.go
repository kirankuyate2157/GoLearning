package routes

import (
	"app/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterTypeRoutes(e *echo.Echo) {
	// Base group: /api/types
	g := e.Group("/api/types")

	// ✅ Single GET endpoint with optional ?category=
	g.GET("", controllers.ListTypes)
	g.POST("", controllers.CreateType)
	g.GET("/:uid", controllers.GetType)
	g.PATCH("/:uid", controllers.UpdateType)
	g.DELETE("/:uid", controllers.DeleteType)

	g.GET("/category", controllers.ListTypeCategory)
}
