package main

import (
	"app/db"
	"app/docs"
	"app/routes"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server for using Swagger with Echo.
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// load .env for local development if present
	if err := godotenv.Load(); err != nil {
		// not fatal; continue and rely on environment
		log.Printf(".env not loaded: %v", err)
	}

	pool := db.DBConfig()
	defer pool.Close()

	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", pool)
			return next(c)
		}
	})
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "ip=${remote_ip}, method=${method}, status=${status}, uri=${uri}, latency=${latency_human}\n",
	}))
	e.Use(middleware.Recover())

	// centralized HTTP error handler to ensure errors are logged
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		// determine status code and message
		code := 500
		var msg interface{}
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			msg = he.Message
		} else if err != nil {
			msg = err.Error()
		}

		// log with both standard log and echo logger
		req := c.Request()
		log.Printf("error response: code=%d method=%s uri=%s remote=%s err=%v", code, req.Method, req.RequestURI, c.RealIP(), msg)
		c.Logger().Errorf("error response: code=%d method=%s uri=%s remote=%s err=%v", code, req.Method, req.RequestURI, c.RealIP(), msg)

		// respond with JSON if possible
		if !c.Response().Committed {
			_ = c.JSON(code, echo.Map{"error": msg})
		}
	}
	swaggerHandler := echoSwagger.WrapHandler
	e.GET("/swagger/*", func(c echo.Context) error {
		docs.SwaggerInfo.Host = c.Request().Host
		docs.SwaggerInfo.Schemes = []string{c.Scheme()} // or []string{"http"} if you only serve http
		return swaggerHandler(c)
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, echo.Map{"message": "pong"})
	})

	routes.VendorRoutes(e)
	routes.RegisterTypeRoutes(e)
	routes.MachineRoutes(e)
	routes.SparePartRoutes(e)
	routes.FeedRoutes(e)
	routes.MaintenanceRoutes(e)
	routes.SettingRoutes(e)
	routes.FileRoutes(e)

	// Serve static files from dist directory
	distDir := os.Getenv("DIST_DIR")
	if distDir == "" {
		// Serve static files
		panic("DIST_DIR env variable not set")
	}
	e.Static("/", distDir)

	// SPA fallback: redirect any non-API requests to index.html
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Call next handler first
			err := next(c)

			// If it's a 404 and not an API request, serve index.html
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok && he.Code == 404 {
					path := c.Request().URL.Path
					// Only redirect non-API routes
					if len(path) < 4 || path[:4] != "/api" {
						return c.File(distDir + "/index.html")
					}
				}
			}
			return err
		}
	})

	e.Logger.Fatal(e.Start(":8080"))
}
