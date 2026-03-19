package service

import (
	"context"
	"errors"
	"springboard/internal/lib"
	"springboard/internal/user/dto"
	"springboard/internal/user/repository"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserService interface {
	GetMe(ctx context.Context, userID string) (dto.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetMe(ctx context.Context, userID string) (dto.UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return dto.UserResponse{}, ErrUserNotFound
	}

	return dto.UserResponse{
		ID:       user.ID,
		Role:     lib.UserRole(user.Role),
		Email:    user.Email,
		FullName: user.FullName,
	}, nil
}
