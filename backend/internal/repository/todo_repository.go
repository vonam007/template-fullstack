package repository

import (
	"context"
	"fmt"
	"time"

	"template-fullstack/backend/internal/domain"
	"template-fullstack/backend/internal/pkg/db"

	"github.com/google/uuid"
)

type TodoRepository interface {
	Create(ctx context.Context, req domain.CreateTodoRequest) (*domain.Todo, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, pagination domain.PaginationQuery) ([]domain.Todo, int64, error)
	Update(ctx context.Context, id uuid.UUID, req domain.UpdateTodoRequest) (*domain.Todo, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, pagination domain.PaginationQuery) ([]domain.Todo, int64, error)
}

type todoRepository struct {
	db *db.DB
}

func NewTodoRepository(database *db.DB) TodoRepository {
	return &todoRepository{db: database}
}

func (r *todoRepository) Create(ctx context.Context, req domain.CreateTodoRequest) (*domain.Todo, error) {
	todo := &domain.Todo{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		UserID:      req.UserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `
		INSERT INTO todos (id, title, description, completed, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, title, description, completed, user_id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, todo.ID, todo.Title, todo.Description, todo.Completed, todo.UserID, todo.CreatedAt, todo.UpdatedAt).
		Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.UserID, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	return todo, nil
}

func (r *todoRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	todo := &domain.Todo{}
	query := `
		SELECT id, title, description, completed, user_id, created_at, updated_at
		FROM todos
		WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.UserID, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get todo by id: %w", err)
	}

	return todo, nil
}

func (r *todoRepository) GetByUserID(ctx context.Context, userID uuid.UUID, pagination domain.PaginationQuery) ([]domain.Todo, int64, error) {
	// Count total records for the user
	var total int64
	countQuery := `SELECT COUNT(*) FROM todos WHERE user_id = $1`
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count todos: %w", err)
	}

	// Get paginated results
	offset := (pagination.Page - 1) * pagination.PageSize
	query := `
		SELECT id, title, description, completed, user_id, created_at, updated_at
		FROM todos
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, userID, pagination.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get todos: %w", err)
	}
	defer rows.Close()

	var todos []domain.Todo
	for rows.Next() {
		var todo domain.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.UserID, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}

	return todos, total, nil
}

func (r *todoRepository) Update(ctx context.Context, id uuid.UUID, req domain.UpdateTodoRequest) (*domain.Todo, error) {
	todo := &domain.Todo{}
	query := `
		UPDATE todos
		SET title = COALESCE(NULLIF($2, ''), title),
		    description = COALESCE(NULLIF($3, ''), description),
		    completed = COALESCE($4, completed),
		    updated_at = $5
		WHERE id = $1
		RETURNING id, title, description, completed, user_id, created_at, updated_at`

	var completed interface{}
	if req.Completed != nil {
		completed = *req.Completed
	}

	err := r.db.QueryRow(ctx, query, id, req.Title, req.Description, completed, time.Now()).
		Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.UserID, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	return todo, nil
}

func (r *todoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM todos WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("todo not found")
	}

	return nil
}

func (r *todoRepository) List(ctx context.Context, pagination domain.PaginationQuery) ([]domain.Todo, int64, error) {
	// Count total records
	var total int64
	countQuery := `SELECT COUNT(*) FROM todos`
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count todos: %w", err)
	}

	// Get paginated results
	offset := (pagination.Page - 1) * pagination.PageSize
	query := `
		SELECT id, title, description, completed, user_id, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, pagination.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get todos: %w", err)
	}
	defer rows.Close()

	var todos []domain.Todo
	for rows.Next() {
		var todo domain.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.UserID, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}

	return todos, total, nil
}
