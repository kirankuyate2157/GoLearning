package routes

import (
	"app/controllers"

	"github.com/labstack/echo/v4"
)

func MachineRoutes(e *echo.Echo) {
	g := e.Group("/api/machines")
	// g.Use(middleware.JWTAuthMiddleware)
	g.GET("", controllers.GetMachines)
	g.POST("", controllers.CreateMachine)
	g.GET("/:uid", controllers.GetMachine)
	g.PATCH("/:uid", controllers.UpdateMachine)
	g.DELETE("/:uid", controllers.DeleteMachine)

	g.GET("/:machine_uid/locations", controllers.ListMachineLocations)
	g.POST("/:machine_uid/locations", controllers.CreateMachineLocation)
	g.PATCH("/:machine_uid/locations/:location_uid", controllers.UpdateMachineLocation)
	g.DELETE("/:machine_uid/locations/:location_uid", controllers.DeleteMachineLocation)

	g.GET("/:machine_uid/warranties", controllers.ListMachineWarranties)
	g.POST("/:machine_uid/warranties", controllers.CreateMachineWarranty)
	g.PATCH("/:machine_uid/warranties/:warranty_uid", controllers.UpdateMachineWarranty)
	g.DELETE("/:machine_uid/warranties/:warranty_uid", controllers.DeleteMachineWarranty)
}
