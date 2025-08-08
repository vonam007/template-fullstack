package service

import (
	"context"
	"math"

	"template-fullstack/backend/internal/domain"
	"template-fullstack/backend/internal/repository"

	"github.com/google/uuid"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Create(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	// In a real application, hash the password here
	return s.userRepo.Create(ctx, req)
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (*domain.User, error) {
	return s.userRepo.Update(ctx, id, req)
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context, pagination domain.PaginationQuery) (*domain.PaginatedResponse, error) {
	users, total, err := s.userRepo.List(ctx, pagination)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	return &domain.PaginatedResponse{
		Data: users,
		Pagination: domain.Pagination{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
