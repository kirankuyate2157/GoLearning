package controllers

import (
	"net/http"

	"app/db/repo"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type governanceController struct{}

func NewGovernanceController() *governanceController { return &governanceController{} }

// CreateOrganizationRequest
type CreateOrganizationRequest struct {
	Name  string `json:"name"`
	Code  string `json:"code,omitempty"`
	Notes string `json:"notes,omitempty"`
}

func (g *governanceController) CreateOrganization(c echo.Context) error {
	var req CreateOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid payload"})
	}
	ctx := c.Request().Context()
	// get db pool from context
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "missing db pool in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db pool in context"})
	}
	gr := repo.NewGovernanceRepo(pool)
	org, err := gr.CreateOrganization(ctx, req.Name, req.Code, req.Notes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, org)
}

func (g *governanceController) ListOrganizations(c echo.Context) error {
	ctx := c.Request().Context()
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "missing db pool in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db pool in context"})
	}
	gr := repo.NewGovernanceRepo(pool)
	orgs, err := gr.ListOrganizations(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, orgs)
}

func (g *governanceController) GetOrganization(c echo.Context) error {
	uid := c.Param("uid")
	ctx := c.Request().Context()
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "missing db pool in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db pool in context"})
	}
	gr := repo.NewGovernanceRepo(pool)
	org, err := gr.GetOrganizationByUID(ctx, uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if org == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "organization not found"})
	}
	return c.JSON(http.StatusOK, org)
}

// Sites
type CreateSiteRequest struct {
	OrganizationUID string `json:"organization_uid"`
	Name            string `json:"name"`
	Code            string `json:"code,omitempty"`
	Location        string `json:"location,omitempty"`
	Notes           string `json:"notes,omitempty"`
}

func (g *governanceController) CreateSite(c echo.Context) error {
	var req CreateSiteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid payload"})
	}
	ctx := c.Request().Context()
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "missing db pool in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db pool in context"})
	}
	gr := repo.NewGovernanceRepo(pool)
	site, err := gr.CreateSite(ctx, req.OrganizationUID, req.Name, req.Code, req.Location, req.Notes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, site)
}

func (g *governanceController) ListSitesByOrganization(c echo.Context) error {
	orgUID := c.Param("org_uid")
	ctx := c.Request().Context()
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "missing db pool in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db pool in context"})
	}
	gr := repo.NewGovernanceRepo(pool)
	sites, err := gr.ListSitesByOrganization(ctx, orgUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, sites)
}

func (g *governanceController) GetSite(c echo.Context) error {
	uid := c.Param("uid")
	ctx := c.Request().Context()
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "missing db pool in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db pool in context"})
	}
	gr := repo.NewGovernanceRepo(pool)
	s, err := gr.GetSiteByUID(ctx, uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if s == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "site not found"})
	}
	return c.JSON(http.StatusOK, s)
}

// Roles and role assignments
type CreateRoleRequest struct {
	Name  string `json:"name"`
	Notes string `json:"notes,omitempty"`
}

func (g *governanceController) CreateRole(c echo.Context) error {
	var req CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid payload"})
	}
	ctx := c.Request().Context()
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "missing db pool in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db pool in context"})
	}
	gr := repo.NewGovernanceRepo(pool)
	role, err := gr.CreateRole(ctx, req.Name, req.Notes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, role)
}

func (g *governanceController) ListRoles(c echo.Context) error {
	ctx := c.Request().Context()
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "missing db pool in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db pool in context"})
	}
	gr := repo.NewGovernanceRepo(pool)
	roles, err := gr.ListRoles(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, roles)
}

// Assign role payload
type AssignRoleRequest struct {
	RoleUID         string  `json:"role_uid"`
	UserUID         string  `json:"user_uid"`
	OrganizationUID *string `json:"organization_uid,omitempty"`
	SiteUID         *string `json:"site_uid,omitempty"`
}

func (g *governanceController) AssignRole(c echo.Context) error {
	var req AssignRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid payload"})
	}
	ctx := c.Request().Context()
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "missing db pool in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db pool in context"})
	}
	gr := repo.NewGovernanceRepo(pool)
	ra, err := gr.AssignRoleToUser(ctx, req.RoleUID, req.UserUID, req.OrganizationUID, req.SiteUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, ra)
}

func (g *governanceController) ListAssignmentsForUser(c echo.Context) error {
	userUID := c.Param("user_uid")
	ctx := c.Request().Context()
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "missing db pool in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db pool in context"})
	}
	gr := repo.NewGovernanceRepo(pool)
	list, err := gr.ListRoleAssignmentsForUser(ctx, userUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

// Simple middleware stub: checks X-User-Role header for required role. This is a lightweight scaffold until full auth is integrated.
func RequireRole(roleName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// read header
			h := c.Request().Header.Get("X-User-Role")
			if h == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing role header"})
			}
			// naive check: comma-separated list
			// in future: resolve user from auth token and check DB role assignments
			for _, v := range splitAndTrim(h) {
				if v == roleName {
					return next(c)
				}
			}
			return c.JSON(http.StatusForbidden, echo.Map{"error": "insufficient role"})
		}
	}
}

func splitAndTrim(s string) []string {
	var out []string
	for _, p := range pgTrimSplit(s, ',') {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func pgTrimSplit(s string, sep rune) []string {
	var cur string
	var out []string
	for _, r := range s {
		if r == sep {
			out = append(out, trimSpaces(cur))
			cur = ""
			continue
		}
		cur += string(r)
	}
	if cur != "" {
		out = append(out, trimSpaces(cur))
	}
	return out
}

func trimSpaces(s string) string {
	return stringTrim(s)
}

// small wrapper to avoid importing strings in multiple files
func stringTrim(s string) string {
	// simple trim: remove leading/trailing spaces
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\n' || s[0] == '\r') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t' || s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
	}
	return s
}
