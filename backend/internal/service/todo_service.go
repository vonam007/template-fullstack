package service

import (
	"context"
	"math"

	"template-fullstack/backend/internal/domain"
	"template-fullstack/backend/internal/repository"

	"github.com/google/uuid"
)

type todoService struct {
	todoRepo repository.TodoRepository
}

func NewTodoService(todoRepo repository.TodoRepository) TodoService {
	return &todoService{todoRepo: todoRepo}
}

func (s *todoService) Create(ctx context.Context, req domain.CreateTodoRequest) (*domain.Todo, error) {
	return s.todoRepo.Create(ctx, req)
}

func (s *todoService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	return s.todoRepo.GetByID(ctx, id)
}

func (s *todoService) GetByUserID(ctx context.Context, userID uuid.UUID, pagination domain.PaginationQuery) (*domain.PaginatedResponse, error) {
	todos, total, err := s.todoRepo.GetByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	return &domain.PaginatedResponse{
		Data: todos,
		Pagination: domain.Pagination{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *todoService) Update(ctx context.Context, id uuid.UUID, req domain.UpdateTodoRequest) (*domain.Todo, error) {
	return s.todoRepo.Update(ctx, id, req)
}

func (s *todoService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.todoRepo.Delete(ctx, id)
}

func (s *todoService) List(ctx context.Context, pagination domain.PaginationQuery) (*domain.PaginatedResponse, error) {
	todos, total, err := s.todoRepo.List(ctx, pagination)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	return &domain.PaginatedResponse{
		Data: todos,
		Pagination: domain.Pagination{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
