package routes

import (
	"app/controllers"
	"github.com/labstack/echo/v4"
)

func QualityRoutes(e *echo.Echo) {
	qualityGroup := e.Group("/api/quality-inspections")
	qualityGroup.GET("", controllers.GetQualityInspections)
	qualityGroup.POST("", controllers.CreateQualityInspection)
	qualityGroup.GET("/:uid", controllers.GetQualityInspection)
	qualityGroup.PATCH("/:uid", controllers.UpdateQualityInspection)
	qualityGroup.DELETE("/:uid", controllers.DeleteQualityInspection)
	qualityGroup.GET("/:uid/non-conformances", controllers.GetQualityInspectionNonConformances)
	qualityGroup.POST("/:uid/non-conformances", controllers.CreateQualityInspectionNonConformance)
	qualityGroup.PATCH("/:uid/non-conformances/:nc_uid", controllers.UpdateQualityInspectionNonConformance)
	qualityGroup.DELETE("/:uid/non-conformances/:nc_uid", controllers.DeleteQualityInspectionNonConformance)
}
