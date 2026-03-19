package service

import (
	"context"
	"errors"
	"log"
	"springboard/internal/lib"
	"springboard/internal/user/repository"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserService interface {
	GetMe(ctx context.Context, userID string, role lib.UserRole) (interface{}, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetMe(ctx context.Context, userID string, role lib.UserRole) (interface{}, error) {
	log.Print("GetMe service start")
	return s.repo.GetFullUserByID(ctx, userID, role)
}
