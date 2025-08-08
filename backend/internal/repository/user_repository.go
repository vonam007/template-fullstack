package repository

import (
	"context"
	"fmt"
	"time"

	"template-fullstack/backend/internal/domain"
	"template-fullstack/backend/internal/pkg/db"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, pagination domain.PaginationQuery) ([]domain.User, int64, error)
}

type userRepository struct {
	db *db.DB
}

func NewUserRepository(database *db.DB) UserRepository {
	return &userRepository{db: database}
}

func (r *userRepository) Create(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	user := &domain.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Name:      req.Name,
		Password:  req.Password, // In real app, this should be hashed
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO users (id, email, name, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, email, name, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, user.ID, user.Email, user.Name, user.Password, user.CreatedAt, user.UpdatedAt).
		Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, email, name, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, email, name, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (*domain.User, error) {
	user := &domain.User{}
	query := `
		UPDATE users
		SET name = COALESCE(NULLIF($2, ''), name),
		    email = COALESCE(NULLIF($3, ''), email),
		    updated_at = $4
		WHERE id = $1
		RETURNING id, email, name, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, id, req.Name, req.Email, time.Now()).
		Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) List(ctx context.Context, pagination domain.PaginationQuery) ([]domain.User, int64, error) {
	// Count total records
	var total int64
	countQuery := `SELECT COUNT(*) FROM users`
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get paginated results
	offset := (pagination.Page - 1) * pagination.PageSize
	query := `
		SELECT id, email, name, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, pagination.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, total, nil
}
