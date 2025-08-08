package middleware

import (
	"net/http"
	"strings"

	"template-fullstack/backend/internal/domain"
	"template-fullstack/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, domain.APIResponse{
				Success: false,
				Error: &domain.APIError{
					Code:    domain.ErrCodeUnauthorized,
					Message: "Authorization header is required",
				},
			})
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, domain.APIResponse{
				Success: false,
				Error: &domain.APIError{
					Code:    domain.ErrCodeUnauthorized,
					Message: "Invalid authorization header format",
				},
			})
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, domain.APIResponse{
				Success: false,
				Error: &domain.APIError{
					Code:    domain.ErrCodeUnauthorized,
					Message: "Token is required",
				},
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, domain.APIResponse{
				Success: false,
				Error: &domain.APIError{
					Code:    domain.ErrCodeUnauthorized,
					Message: "Invalid token",
				},
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
