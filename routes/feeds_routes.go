package routes

import (
	"app/controllers"

	"github.com/labstack/echo/v4"
)

func FeedRoutes(e *echo.Echo) {
	// Base group: /api/Feeds

	g := e.Group("/api/feeds")

	g.GET("", controllers.ListFeeds)
	g.POST("", controllers.CreateFeed)
	g.GET("/:uid", controllers.GetFeed)
	g.PATCH("/:uid", controllers.UpdateFeed)
	g.DELETE("/:uid", controllers.DeleteFeed)

	// Feed Machine Mapping Routes
	g.GET("/:feed_uid/machines", controllers.ListFeedMachine)
	g.GET("/machines/:machine_uid", controllers.ListMachineFeed)
	g.POST("/:feed_uid/machines/:machine_uid", controllers.CreateFeedMachine)
	g.DELETE("/:feed_uid/machines/:machine_uid", controllers.RemoveFeedFromMachine)

	g.GET("/type", controllers.ListFeedTypes)
	g.POST("/type", controllers.CreateFeedType)
	g.GET("/type/:uid", controllers.GetFeedType)
	g.PATCH("/type/:uid", controllers.UpdateFeedType)
	g.DELETE("/type/:uid", controllers.DeleteFeedType)

	g.GET("/alert", controllers.ListFeedsAlert)
	g.POST("/:feed_uid/alert", controllers.CreateFeedsAlert)
	g.GET("/:feed_uid/alert", controllers.GetFeedsAlert)
	g.PATCH("/:feed_uid/alert/:uid", controllers.UpdateFeedsAlert)
	g.DELETE("/alert/:uid", controllers.DeleteFeedsAlert)

	// Feed Ingestion Routes
	g.GET("/ingestions", controllers.ListActiveIngestions)
	g.POST("/:uid/ingestions", controllers.StartIngestion)
	g.DELETE("/:uid/ingestions", controllers.StopIngestion)
}
