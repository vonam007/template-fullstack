package service

import (
	"context"
	"fmt"
	"time"

	"template-fullstack/backend/internal/domain"
	"template-fullstack/backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error)
	ValidateToken(tokenString string) (*Claims, error)
	GenerateToken(userID uuid.UUID, email string) (string, error)
}

type UserService interface {
	Create(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Update(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, pagination domain.PaginationQuery) (*domain.PaginatedResponse, error)
}

type TodoService interface {
	Create(ctx context.Context, req domain.CreateTodoRequest) (*domain.Todo, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, pagination domain.PaginationQuery) (*domain.PaginatedResponse, error)
	Update(ctx context.Context, id uuid.UUID, req domain.UpdateTodoRequest) (*domain.Todo, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, pagination domain.PaginationQuery) (*domain.PaginatedResponse, error)
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpiry time.Duration) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	// In a real application, you would hash the password and compare
	// For this template, we'll do a simple mock validation
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Simple password check (in real app, use bcrypt)
	if user.Password != req.Password {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := s.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) GenerateToken(userID uuid.UUID, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
