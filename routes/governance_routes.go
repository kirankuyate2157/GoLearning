package routes

import (
	"app/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterGovernanceRoutes(e *echo.Echo) {
	ctl := controllers.NewGovernanceController()
	g := e.Group("/api/governance")

	g.POST("/organizations", ctl.CreateOrganization)
	g.GET("/organizations", ctl.ListOrganizations)
	g.GET("/organizations/:uid", ctl.GetOrganization)

	g.POST("/sites", ctl.CreateSite)
	g.GET("/organizations/:org_uid/sites", ctl.ListSitesByOrganization)
	g.GET("/sites/:uid", ctl.GetSite)

	g.POST("/roles", ctl.CreateRole)
	g.GET("/roles", ctl.ListRoles)

	g.POST("/role-assignments", ctl.AssignRole)
	g.GET("/role-assignments/user/:user_uid", ctl.ListAssignmentsForUser)
}
