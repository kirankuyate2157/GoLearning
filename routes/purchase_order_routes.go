package routes

import (
	"app/controllers"

	"github.com/labstack/echo/v4"
)

func PurchaseOrderRoutes(e *echo.Echo) {
	purchaseGroup := e.Group("/api/purchase-orders")
	purchaseGroup.GET("", controllers.GetPurchaseOrders)
	purchaseGroup.POST("", controllers.CreatePurchaseOrder)
	purchaseGroup.GET("/:uid", controllers.GetPurchaseOrder)
	purchaseGroup.PATCH("/:uid", controllers.UpdatePurchaseOrder)
	purchaseGroup.DELETE("/:uid", controllers.DeletePurchaseOrder)
	purchaseGroup.GET("/:uid/lines", controllers.GetPurchaseOrderLines)
	purchaseGroup.POST("/:uid/lines", controllers.CreatePurchaseOrderLine)
	purchaseGroup.PATCH("/:uid/lines/:line_uid", controllers.UpdatePurchaseOrderLine)
	purchaseGroup.DELETE("/:uid/lines/:line_uid", controllers.DeletePurchaseOrderLine)

	vendorGroup := e.Group("/api/vendors")
	vendorGroup.GET("/:vendor_uid/purchase-orders", controllers.GetVendorPurchaseOrders)
}
