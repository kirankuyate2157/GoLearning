package routes

import (
	"app/controllers"

	"github.com/labstack/echo/v4"
)

func VendorRoutes(e *echo.Echo) {
	g := e.Group("/api/vendors")
	// g.Use(middleware.JWTAuthMiddleware)
	g.GET("", controllers.GetVendors)
	g.POST("", controllers.CreateVendor)
	g.GET(":uid", controllers.GetVendor)
	g.PATCH(":uid", controllers.UpdateVendor)
	g.DELETE(":uid", controllers.DeleteVendor)

	g.GET("/:vendor_uid/addresses", controllers.ListVendorAddresses)
	g.POST("/:vendor_uid/addresses", controllers.CreateVendorAddress)
	// g.GET(":id", controllers.GetVendorAddress)
	g.PATCH("/:vendor_uid/addresses/:address_uid", controllers.UpdateVendorAddress)
	g.DELETE("/:vendor_uid/addresses/:address_uid", controllers.DeleteVendorAddress)

	// Vendor contact routes
	g.GET("/:vendor_uid/contacts", controllers.ListVendorContacts)
	g.POST("/:vendor_uid/contacts", controllers.CreateVendorContact)
	g.PATCH("/:vendor_uid/contacts/:contact_uid", controllers.UpdateVendorContact)
	g.DELETE("/:vendor_uid/contacts/:contact_uid", controllers.DeleteVendorContact)

	// Vendor Type routes
	g.GET("/:vendor_uid/vendorsTypes", controllers.GetVendorTypes)
	g.POST("/:vendor_uid/vendorsTypes", controllers.CreateVendorType)
	g.DELETE("/:vendor_uid/vendorsTypes/:vendorAssign_uid", controllers.DeleteVendorType)

}
