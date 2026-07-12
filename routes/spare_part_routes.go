package routes

import (
	"app/controllers"

	"github.com/labstack/echo/v4"
)

func SparePartRoutes(e *echo.Echo) {
	g := e.Group("/api/spare-parts")

	// Spare Part
	g.GET("", controllers.GetSpareParts)
	g.POST("", controllers.CreateSparePart)
	g.GET("/:uid", controllers.GetSparePart)
	g.PATCH("/:uid", controllers.UpdateSparePart)
	g.DELETE("/:uid", controllers.DeleteSparePart)

	// Spare Part Vendor
	g.GET("/:sparepart_uid/vendors", controllers.GetSparePartVendorList)
	g.POST("/:sparepart_uid/vendors", controllers.CreateSparePartVendor)
	g.PATCH("/:sparepart_uid/vendors/:sparepartvendor_uid", controllers.UpdateSparePartVendor)
	g.DELETE("/:sparepart_uid/vendors/:sparepartvendor_uid", controllers.DeleteSparePartVendor)

	//Spare part machine
	g.GET("/:sparepart_uid/machines", controllers.GetSparePartMachineList)
	g.POST("/:sparepart_uid/machines", controllers.CreateSparePartMachine)
	g.DELETE("/:sparepart_uid/machines/:sparepartmachine_uid", controllers.DeleteSparePartMachine)

}
