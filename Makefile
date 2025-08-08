.PHONY: help dev prod build clean test lint migrate-up migrate-down seed logs

# Default target
help: ## Display this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development
dev: ## Start development environment
	@echo "Starting development environment..."
	@docker compose --profile dev up --build -d
	@echo "✅ Development environment started!"
	@echo "Frontend: http://localhost:5173"
	@echo "Backend: http://localhost:8080"
	@echo "API Docs: http://localhost:8080/swagger/index.html"

prod: ## Start production environment
	@echo "Starting production environment..."
	@docker compose --profile prod up --build -d
	@echo "✅ Production environment started!"
	@echo "Application: http://localhost"
	@echo "API: http://localhost/api"

# Backend commands
be.run: ## Run backend locally
	@cd backend && go run cmd/server/main.go

be.build: ## Build backend
	@cd backend && go build -o bin/server cmd/server/main.go

be.test: ## Run backend tests
	@cd backend && go test -v ./...

be.lint: ## Lint backend code
	@cd backend && golangci-lint run

# Frontend commands
fe.dev: ## Run frontend in development mode
	@cd frontend && npm run dev

fe.build: ## Build frontend for production
	@cd frontend && npm run build

fe.test: ## Run frontend tests
	@cd frontend && npm test

fe.lint: ## Lint frontend code
	@cd frontend && npm run lint

# Database commands
db.migrate.up: ## Run database migrations up
	@echo "Running database migrations..."
	@docker compose exec backend migrate -path=/app/migrations -database="postgres://$(shell grep DB_USER .env | cut -d '=' -f2):$(shell grep DB_PASSWORD .env | cut -d '=' -f2)@db:5432/$(shell grep DB_NAME .env | cut -d '=' -f2)?sslmode=disable" up
	@echo "✅ Migrations completed!"

db.migrate.down: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	@docker compose exec backend migrate -path=/app/migrations -database="postgres://$(shell grep DB_USER .env | cut -d '=' -f2):$(shell grep DB_PASSWORD .env | cut -d '=' -f2)@db:5432/$(shell grep DB_NAME .env | cut -d '=' -f2)?sslmode=disable" down 1
	@echo "✅ Rollback completed!"

db.seed: ## Seed database with sample data
	@echo "Seeding database..."
	@docker compose exec backend go run cmd/seed/main.go
	@echo "✅ Database seeded!"

db.reset: ## Reset database (down all migrations and up again)
	@echo "Resetting database..."
	@docker compose exec backend migrate -path=/app/migrations -database="postgres://$(shell grep DB_USER .env | cut -d '=' -f2):$(shell grep DB_PASSWORD .env | cut -d '=' -f2)@db:5432/$(shell grep DB_NAME .env | cut -d '=' -f2)?sslmode=disable" down -all
	@docker compose exec backend migrate -path=/app/migrations -database="postgres://$(shell grep DB_USER .env | cut -d '=' -f2):$(shell grep DB_PASSWORD .env | cut -d '=' -f2)@db:5432/$(shell grep DB_NAME .env | cut -d '=' -f2)?sslmode=disable" up
	@echo "✅ Database reset completed!"

# Utility commands
logs: ## Show logs for all services
	@docker compose logs -f

logs.backend: ## Show backend logs
	@docker compose logs -f backend

logs.frontend: ## Show frontend logs
	@docker compose logs -f frontend-dev

logs.db: ## Show database logs
	@docker compose logs -f db

clean: ## Clean up containers and volumes
	@echo "Cleaning up..."
	@docker compose down -v
	@docker system prune -f
	@echo "✅ Cleanup completed!"

build: ## Build all services
	@echo "Building all services..."
	@docker compose build
	@echo "✅ Build completed!"

test: ## Run all tests
	@echo "Running all tests..."
	@make be.test
	@make fe.test
	@echo "✅ All tests completed!"

lint: ## Run linting for all services
	@echo "Running linting..."
	@make be.lint
	@make fe.lint
	@echo "✅ Linting completed!"

# Setup commands
setup: ## Initial setup (copy .env, install dependencies)
	@echo "Setting up project..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "📄 Created .env file"; fi
	@cd frontend && npm install
	@echo "✅ Setup completed!"
	@echo "Next steps:"
	@echo "1. Edit .env file with your configuration"
	@echo "2. Run 'make dev' to start development environment"
	@echo "3. Run 'make db.migrate.up' to setup database"
	@echo "4. Run 'make db.seed' to add sample data"
