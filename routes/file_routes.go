package routes

import (
	"app/controllers"

	"github.com/labstack/echo/v4"
)

func FileRoutes(e *echo.Echo) {
	g := e.Group("/api/files")

	g.POST("", controllers.UploadFile)
	g.PATCH("/:uid", controllers.PatchFile)
	g.GET("/:uid", controllers.ViewFile)
	g.DELETE("/:uid", controllers.DeleteFile)
}
