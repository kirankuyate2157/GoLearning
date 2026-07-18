package routes

import (
	"app/controllers"

	"github.com/labstack/echo/v4"
)

func InventoryRoutes(e *echo.Echo) {
	inventoryGroup := e.Group("/api/inventory")
	inventoryGroup.GET("", controllers.GetInventoryItems)
	inventoryGroup.POST("", controllers.CreateInventoryItem)
	inventoryGroup.GET("/:uid", controllers.GetInventoryItem)
	inventoryGroup.PATCH("/:uid", controllers.UpdateInventoryItem)
	inventoryGroup.DELETE("/:uid", controllers.DeleteInventoryItem)

	inventoryGroup.GET("/:inventory_uid/movements", controllers.GetInventoryMovements)
	inventoryGroup.POST("/:inventory_uid/movements", controllers.CreateInventoryMovement)

	wipGroup := e.Group("/api/wip")
	wipGroup.GET("", controllers.GetWIPItems)
	wipGroup.POST("", controllers.CreateWIPItem)
	wipGroup.GET("/:uid", controllers.GetWIPItem)
	wipGroup.PATCH("/:uid", controllers.UpdateWIPItem)
	wipGroup.DELETE("/:uid", controllers.DeleteWIPItem)
}
