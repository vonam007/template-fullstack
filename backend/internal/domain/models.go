package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	Password  string    `json:"-" db:"password_hash"` // Don't expose password in JSON
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Todo represents a todo item
type Todo struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Completed   bool      `json:"completed" db:"completed"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// DTOs for requests and responses

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=2"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" binding:"omitempty,min=2"`
	Email string `json:"email,omitempty" binding:"omitempty,email"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateTodoRequest struct {
	Title       string    `json:"title" binding:"required,min=1"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"user_id"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title,omitempty" binding:"omitempty,min=1"`
	Description string `json:"description,omitempty"`
	Completed   *bool  `json:"completed,omitempty"`
}

// API Response wrappers
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Pagination
type PaginationQuery struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=10" binding:"min=1,max=100"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Constants for error codes
const (
	ErrCodeInvalidRequest     = "INVALID_REQUEST"
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeForbidden          = "FORBIDDEN"
	ErrCodeNotFound           = "NOT_FOUND"
	ErrCodeConflict           = "CONFLICT"
	ErrCodeInternalError      = "INTERNAL_ERROR"
	ErrCodeTodoNotFound       = "TODO_NOT_FOUND"
	ErrCodeUserNotFound       = "USER_NOT_FOUND"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeEmailExists        = "EMAIL_EXISTS"
)
