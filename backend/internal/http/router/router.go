package router

import (
	"net/http"
	"time"

	"template-fullstack/backend/internal/config"
	"template-fullstack/backend/internal/http/handlers"
	"template-fullstack/backend/internal/http/middleware"
	"template-fullstack/backend/internal/pkg/db"
	"template-fullstack/backend/internal/repository"
	"template-fullstack/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New(cfg *config.Config, database *db.DB, log zerolog.Logger) *gin.Engine {
	// Set Gin mode
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Middleware
	r.Use(gin.Recovery())
	r.Use(middleware.StructuredLoggingMiddleware(log))
	r.Use(middleware.ErrorHandlingMiddleware())
	r.Use(middleware.CORSMiddleware(cfg.CORS.Origins))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	// Initialize repositories
	userRepo := repository.NewUserRepository(database)
	todoRepo := repository.NewTodoRepository(database)

	// Initialize services
	jwtExpiry, _ := time.ParseDuration(cfg.JWT.Expiry)
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret, jwtExpiry)
	// userService := service.NewUserService(userRepo) // TODO: Add user management handlers
	todoService := service.NewTodoService(todoRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	todoHandler := handlers.NewTodoHandler(todoService)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api")
	v1 := api.Group("/v1")

	// Auth routes (public)
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", middleware.AuthMiddleware(authService), authHandler.Refresh)
	}

	// Todo routes (protected)
	todos := v1.Group("/todos")
	todos.Use(middleware.AuthMiddleware(authService))
	{
		todos.POST("", todoHandler.CreateTodo)
		todos.GET("", todoHandler.GetTodos)
		todos.GET("/:id", todoHandler.GetTodo)
		todos.PUT("/:id", todoHandler.UpdateTodo)
		todos.DELETE("/:id", todoHandler.DeleteTodo)
	}

	// Admin routes (protected)
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthMiddleware(authService))
	{
		admin.GET("/todos", todoHandler.GetAllTodos)
	}

	return r
}
