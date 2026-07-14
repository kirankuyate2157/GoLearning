package routes

import (
	"app/controllers"

	"github.com/labstack/echo/v4"
)

func ProductionRoutes(e *echo.Echo) {
	productionGroup := e.Group("/api/production-orders")
	productionGroup.GET("", controllers.GetProductionOrders)
	productionGroup.POST("", controllers.CreateProductionOrder)
	productionGroup.GET("/:uid", controllers.GetProductionOrder)
	productionGroup.PATCH("/:uid", controllers.UpdateProductionOrder)
	productionGroup.DELETE("/:uid", controllers.DeleteProductionOrder)
	productionGroup.GET("/:uid/jobs", controllers.GetProductionOrderJobs)

	jobsGroup := e.Group("/api/jobs")
	jobsGroup.GET("", controllers.GetJobs)
	jobsGroup.POST("", controllers.CreateJob)
	jobsGroup.GET("/:uid", controllers.GetJob)
	jobsGroup.PATCH("/:uid", controllers.UpdateJob)
	jobsGroup.DELETE("/:uid", controllers.DeleteJob)
}
