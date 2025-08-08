package handlers

import (
	"net/http"

	"template-fullstack/backend/internal/domain"
	"template-fullstack/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login credentials"
// @Success 200 {object} domain.APIResponse{data=domain.LoginResponse}
// @Failure 400 {object} domain.APIResponse{error=domain.APIError}
// @Failure 401 {object} domain.APIResponse{error=domain.APIError}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeInvalidRequest,
				Message: err.Error(),
			},
		})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeInvalidCredentials,
				Message: "Invalid email or password",
			},
		})
		return
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success: true,
		Data:    resp,
	})
}

// Mock refresh endpoint
// @Summary Refresh token
// @Description Refresh JWT token
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.APIResponse{data=map[string]string}
// @Failure 401 {object} domain.APIResponse{error=domain.APIError}
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	// Mock implementation - in real app, validate refresh token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeUnauthorized,
				Message: "Unauthorized",
			},
		})
		return
	}

	email, _ := c.Get("email")

	// Convert userID to UUID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeInvalidRequest,
				Message: "Invalid user ID",
			},
		})
		return
	}

	// Generate new token
	token, err := h.authService.GenerateToken(userUUID, email.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIResponse{
			Success: false,
			Error: &domain.APIError{
				Code:    domain.ErrCodeInternalError,
				Message: "Failed to generate token",
			},
		})
		return
	}

	c.JSON(http.StatusOK, domain.APIResponse{
		Success: true,
		Data:    map[string]string{"token": token},
	})
}
