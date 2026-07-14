package routes

import (
	"app/controllers"

	"github.com/labstack/echo/v4"
)

func BOMRoutes(e *echo.Echo) {
	bomsGroup := e.Group("/api/boms")
	bomsGroup.GET("", controllers.GetBOMs)
	bomsGroup.POST("", controllers.CreateBOM)
	bomsGroup.GET("/:uid", controllers.GetBOM)
	bomsGroup.PATCH("/:uid", controllers.UpdateBOM)
	bomsGroup.DELETE("/:uid", controllers.DeleteBOM)

	bomsGroup.GET("/:bom_uid/lines", controllers.GetBOMLines)
	bomsGroup.POST("/:bom_uid/lines", controllers.CreateBOMLine)
	bomsGroup.GET("/:bom_uid/lines/:line_uid", controllers.GetBOMLine)
	bomsGroup.PATCH("/:bom_uid/lines/:line_uid", controllers.UpdateBOMLine)
	bomsGroup.DELETE("/:bom_uid/lines/:line_uid", controllers.DeleteBOMLine)

	routingGroup := e.Group("/api/routing-operations")
	routingGroup.GET("", controllers.GetRoutingOperations)
	routingGroup.POST("", controllers.CreateRoutingOperation)
	routingGroup.GET("/:uid", controllers.GetRoutingOperation)
	routingGroup.PATCH("/:uid", controllers.UpdateRoutingOperation)
	routingGroup.DELETE("/:uid", controllers.DeleteRoutingOperation)
}
