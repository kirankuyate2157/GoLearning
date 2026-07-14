package controllers

import (
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"
)

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns service health and uptime status
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"status":  "ok",
		"service": "industrial-automation-backend",
		"time":    time.Now().UTC(),
	})
}

// ReadinessCheck godoc
// @Summary Readiness check endpoint
// @Description Returns readiness status for dependent services
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health/ready [get]
func ReadinessCheck(c echo.Context) error {
	ready := true
	if c.Get("db") == nil {
		ready = false
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status":    "ready",
		"database":  ready,
		"timestamp": time.Now().UTC(),
	})
}

// VersionInfo godoc
// @Summary Service version info
// @Description Returns backend service version and runtime metadata
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health/version [get]
func VersionInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"service":    "industrial-automation-backend",
		"go_version": runtime.Version(),
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
		"build_time": os.Getenv("BUILD_TIME"),
		"release":    os.Getenv("RELEASE_TAG"),
	})
}
