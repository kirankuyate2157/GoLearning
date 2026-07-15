package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"app/models"
	"github.com/labstack/echo/v4"
)

// UserContext holds the resolved user identity from JWT or session
type UserContext struct {
	UserUID         string
	Username        string
	Email           string
	Roles           []string
	OrganizationUID string
}

// JWTAuthMiddleware extracts and validates JWT token, populates user context in request
// Token is expected in Authorization header: "Bearer <token>"
// For now this uses a simple JWT parsing (no signature validation) as a scaffold
// In production: integrate with oauth2, use proper key management, etc
func JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			// allow unauthenticated requests; controllers can check context
			return next(c)
		}

		// parse "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid auth header format"})
		}

		tokenString := parts[1]

		// simple JWT parsing scaffold: decode without verification (production should verify signature)
		userCtx, err := parseJWTToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid or expired token"})
		}

		// populate context for downstream handlers
		c.Set("user", userCtx)
		return next(c)
	}
}

// parseJWTToken decodes a JWT token and extracts claims (scaffold: no signature verification)
// in production: use proper JWT validation with HMAC/RSA keys
func parseJWTToken(tokenString string) (*UserContext, error) {
	// For this scaffold, we do basic parsing
	// Production should use github.com/golang-jwt/jwt or similar
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid token format")
	}

	// Placeholder: parse claims from payload
	// In production, decode base64, verify signature, extract claims
	var claims models.JWTClaims
	_ = json.Unmarshal([]byte("{}"), &claims) // scaffold: dummy parse

	// For now, return a basic context
	// When integrated with a real auth provider, this will populate from actual JWT
	return &UserContext{
		UserUID:         claims.UserUID,
		Username:        claims.Username,
		Email:           claims.Email,
		Roles:           claims.Roles,
		OrganizationUID: claims.OrganizationUID,
	}, nil
}

// GetUserContext retrieves user context from request. Returns nil if not authenticated
func GetUserContext(c echo.Context) *UserContext {
	if u := c.Get("user"); u != nil {
		if uc, ok := u.(*UserContext); ok {
			return uc
		}
	}
	return nil
}

// RequireAuthMiddleware ensures user is authenticated (user context exists)
func RequireAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if GetUserContext(c) == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "authentication required"})
		}
		return next(c)
	}
}

// RequireRoleMiddleware checks if user has a required role
func RequireRoleMiddleware(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := GetUserContext(c)
			if user == nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "authentication required"})
			}

			// check if user has required role
			for _, r := range user.Roles {
				if r == requiredRole {
					return next(c)
				}
			}
			return c.JSON(http.StatusForbidden, echo.Map{"error": "insufficient role"})
		}
	}
}
