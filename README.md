# Fullstack Template

A production-ready monorepo template featuring:

- **Backend**: Golang (Gin) with REST API, Swagger documentation, JWT authentication, PostgreSQL
- **Frontend**: React + TypeScript + Vite with Redux Toolkit, React Router, i18next
- **Database**: PostgreSQL with migrations
- **Infrastructure**: Docker Compose, Nginx reverse proxy
- **Development**: Hot reload, structured logging, error handling

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Nginx (80)    │────│  Backend (8080) │────│ PostgreSQL      │
│                 │    │                 │    │ (5432)          │
│ Reverse Proxy   │    │ Gin + REST API  │    │                 │
│ Static Files    │    │ JWT Auth        │    │ Migrations      │
└─────────────────┘    │ Swagger Docs    │    │ Seed Data       │
         │              └─────────────────┘    └─────────────────┘
         │              
┌─────────────────┐    
│ Frontend (5173) │    
│                 │    
│ React + TypeScript │    
│ Redux Toolkit   │    
│ React Router    │    
│ i18next         │    
└─────────────────┘    
```

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Make (optional, for convenience commands)

### 1. Clone and Setup

```bash
git clone <repository-url>
cd template-fullstack
cp .env.example .env
```

### 2. Development Environment

```bash
# Using Make (recommended)
make setup
make dev

# Or using Docker Compose directly
docker compose --profile dev up --build
```

### 3. Database Setup

```bash
# Run migrations
make db.migrate.up

# Seed with sample data
make db.seed
```

### 4. Access the Application

- **Frontend**: http://localhost (via Nginx) or http://localhost:5173 (direct)
- **Backend API**: http://localhost/api/v1
- **API Documentation**: http://localhost/swagger/index.html
- **Database**: localhost:5432

### Demo Credentials

- Email: `admin@example.com`
- Password: `admin123`

## 🛠️ Development

### Available Commands

```bash
# Development
make dev              # Start development environment
make prod             # Start production environment

# Backend
make be.run           # Run backend locally
make be.build         # Build backend
make be.test          # Run backend tests
make be.lint          # Lint backend code

# Frontend
make fe.dev           # Run frontend development server
make fe.build         # Build frontend for production
make fe.test          # Run frontend tests
make fe.lint          # Lint frontend code

# Database
make db.migrate.up    # Run database migrations
make db.migrate.down  # Rollback last migration
make db.seed          # Seed database with sample data
make db.reset         # Reset database (down all + up)

# Utilities
make logs             # Show logs for all services
make logs.backend     # Show backend logs only
make logs.frontend    # Show frontend logs only
make clean            # Clean up containers and volumes
make build            # Build all services
make test             # Run all tests
```

### Project Structure

```
.
├── backend/                 # Golang backend
│   ├── cmd/
│   │   ├── server/         # Main application
│   │   └── seed/           # Database seeding
│   ├── internal/
│   │   ├── config/         # Configuration management
│   │   ├── domain/         # Domain models and DTOs
│   │   ├── http/
│   │   │   ├── handlers/   # HTTP handlers
│   │   │   ├── middleware/ # HTTP middleware
│   │   │   └── router/     # Route definitions
│   │   ├── pkg/
│   │   │   ├── db/         # Database connection
│   │   │   └── logger/     # Structured logging
│   │   ├── repository/     # Data access layer
│   │   └── service/        # Business logic layer
│   ├── migrations/         # Database migrations
│   ├── docs/              # Swagger documentation
│   ├── go.mod
│   └── Dockerfile
├── frontend/               # React frontend
│   ├── src/
│   │   ├── app/           # Store, axios, i18n configuration
│   │   ├── components/    # Reusable UI components
│   │   ├── features/      # Feature-based Redux slices
│   │   ├── locales/       # Translation files
│   │   ├── pages/         # Page components
│   │   └── main.tsx
│   ├── package.json
│   ├── vite.config.ts
│   ├── Dockerfile
│   └── Dockerfile.dev
├── deploy/                 # Deployment configuration
│   ├── nginx.conf         # Production Nginx config
│   ├── nginx.dev.conf     # Development Nginx config
│   ├── Dockerfile.nginx
│   └── Dockerfile.nginx.dev
├── docker-compose.yml     # Multi-service orchestration
├── Makefile              # Development commands
├── .env.example          # Environment variables template
└── README.md
```

## 🔧 Configuration

### Environment Variables

Copy `.env.example` to `.env` and customize:

```env
# Application
APP_ENV=development
APP_PORT=8080
APP_NAME=fullstack-template

# Database
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=fullstack_db

# Authentication
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRY=24h

# CORS
CORS_ORIGINS=http://localhost:3000,http://localhost:5173,http://localhost

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## 🏭 Production Deployment

### 1. Production Build

```bash
make prod
```

This will:
- Build optimized frontend assets
- Create production Docker images
- Start services with Nginx serving static files
- Run database migrations automatically

### 2. Manual Production Steps

```bash
# Build and start production services
docker compose --profile prod up --build -d

# Run database migrations
make db.migrate.up

# Seed initial data (optional)
make db.seed
```

### 3. Health Checks

All services include health checks:
- **Backend**: `GET /health`
- **Frontend/Nginx**: `GET /`
- **Database**: `pg_isready`

### 4. Environment-Specific Configurations

- **Development**: Hot reload, detailed logs, CORS enabled
- **Production**: Optimized builds, compressed assets, security headers

## 🔒 Security Features

- JWT-based authentication
- CORS configuration
- Security headers in Nginx
- Non-root container users
- Input validation
- SQL injection prevention (parameterized queries)

## 📊 API Documentation

Interactive API documentation is available at:
- Development: http://localhost:8080/swagger/index.html
- Production: http://localhost/swagger/index.html

### API Endpoints

#### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh JWT token

#### Todos
- `GET /api/v1/todos` - Get user's todos (paginated)
- `POST /api/v1/todos` - Create new todo
- `GET /api/v1/todos/{id}` - Get specific todo
- `PUT /api/v1/todos/{id}` - Update todo
- `DELETE /api/v1/todos/{id}` - Delete todo

#### Admin
- `GET /api/v1/admin/todos` - Get all todos (admin)

## 🌍 Internationalization

The frontend supports multiple languages:
- English (en)
- Vietnamese (vi)

To add a new language:
1. Create `src/locales/{language}.json`
2. Add translations for all keys
3. Update `LanguageSelector` component

## 🧪 Testing

### Backend Tests
```bash
make be.test
```

### Frontend Tests
```bash
make fe.test
```

### Integration Tests
```bash
# Start test environment
docker compose --profile dev up -d

# Run API tests (example with curl)
curl -X POST http://localhost/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'
```

## 🔍 Monitoring and Logging

### Structured Logging
- All services use structured JSON logging
- Configurable log levels
- Request/response logging
- Error tracking

### Log Access
```bash
# All services
make logs

# Specific service
make logs.backend
make logs.frontend
make logs.db
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure ports 80, 5173, 8080, 5432 are available
2. **Database connection**: Check if PostgreSQL container is running
3. **Migration errors**: Ensure database is accessible before running migrations
4. **Frontend build fails**: Clear node_modules and reinstall dependencies

### Reset Everything

```bash
make clean
docker system prune -f
make setup
make dev
```

### Check Service Status

```bash
docker compose ps
docker compose logs <service-name>
```

## 🔗 Related Documentation

- [Gin Framework](https://gin-gonic.com/)
- [React Documentation](https://reactjs.org/)
- [Redux Toolkit](https://redux-toolkit.js.org/)
- [Vite](https://vitejs.dev/)
- [Docker Compose](https://docs.docker.com/compose/)
- [PostgreSQL](https://www.postgresql.org/docs/)
