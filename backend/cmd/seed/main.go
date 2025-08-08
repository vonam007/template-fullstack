package main

import (
	"context"

	"template-fullstack/backend/internal/config"
	"template-fullstack/backend/internal/domain"
	"template-fullstack/backend/internal/pkg/db"
	"template-fullstack/backend/internal/pkg/logger"
	"template-fullstack/backend/internal/repository"

	"github.com/google/uuid"
)

func main() {
	log := logger.New()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize database
	database, err := db.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(database)
	todoRepo := repository.NewTodoRepository(database)

	ctx := context.Background()

	// Seed users
	users := []domain.CreateUserRequest{
		{
			Email:    "admin@example.com",
			Name:     "Admin User",
			Password: "admin123",
		},
		{
			Email:    "user@example.com",
			Name:     "Regular User",
			Password: "user123",
		},
	}

	var userIDs []uuid.UUID
	for _, userReq := range users {
		user, err := userRepo.Create(ctx, userReq)
		if err != nil {
			log.Error().Err(err).Str("email", userReq.Email).Msg("Failed to create user")
			continue
		}
		userIDs = append(userIDs, user.ID)
		log.Info().Str("email", user.Email).Msg("Created user")
	}

	// Seed todos
	if len(userIDs) > 0 {
		todos := []domain.CreateTodoRequest{
			{
				Title:       "Learn Golang",
				Description: "Study Golang fundamentals and best practices",
				UserID:      userIDs[0],
			},
			{
				Title:       "Build REST API",
				Description: "Create a REST API using Gin framework",
				UserID:      userIDs[0],
			},
			{
				Title:       "Setup Docker",
				Description: "Configure Docker containers for development",
				UserID:      userIDs[1],
			},
			{
				Title:       "Learn React",
				Description: "Study React hooks and Redux Toolkit",
				UserID:      userIDs[1],
			},
		}

		for _, todoReq := range todos {
			todo, err := todoRepo.Create(ctx, todoReq)
			if err != nil {
				log.Error().Err(err).Str("title", todoReq.Title).Msg("Failed to create todo")
				continue
			}
			log.Info().Str("title", todo.Title).Msg("Created todo")
		}
	}

	log.Info().Msg("Database seeding completed!")
}
