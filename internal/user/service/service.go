package service

import (
	"context"
	"errors"
	"log"
	"springboard/internal/lib"
	"springboard/internal/user/dto"
	"springboard/internal/user/repository"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserService interface {
	GetMe(ctx context.Context, userID string, role lib.UserRole) (any, error)
	DeleteMe(ctx context.Context, userID string) error
	UpdateCandidate(ctx context.Context, userID string, data dto.UpdateMeCandidateRequest) (any, error)
	UpdateEmployer(ctx context.Context, userID string, data dto.UpdateMeEmployerRequest) (any, error)
	Verify(ctx context.Context, userID string, role lib.UserRole, inn string) error
	SetPrivacy(ctx context.Context, userID string, isPrivate bool) error
	SetAvatar(ctx context.Context, userID string, role lib.UserRole, url string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetMe(ctx context.Context, userID string, role lib.UserRole) (any, error) {
	log.Print("GetMe service start")
	return s.repo.GetFullUserByID(ctx, userID, role)
}

func (s *userService) DeleteMe(ctx context.Context, userID string) error {
	log.Print(" DeleteMe service start")
	return s.repo.DeleteUserByID(ctx, userID)
}

func (s *userService) UpdateCandidate(ctx context.Context, userID string, data dto.UpdateMeCandidateRequest) (any, error) {
	log.Print("UpdateCandidate service start")

	err := s.repo.UpdateCandidate(ctx, userID, data)
	if err != nil {
		return nil, err
	}
	return s.repo.GetFullUserByID(ctx, userID, lib.RoleStudent)
}

func (s *userService) UpdateEmployer(ctx context.Context, userID string, data dto.UpdateMeEmployerRequest) (any, error) {
	log.Print("UpdateEmployer service start")

	err := s.repo.UpdateEmployer(ctx, userID, data)
	if err != nil {
		return nil, err
	}

	return s.repo.GetFullUserByID(ctx, userID, lib.RoleEmployer)
}

func (s *userService) Verify(ctx context.Context, userID string, role lib.UserRole, inn string) error {
	if role != lib.RoleEmployer {
		return errors.New("only employers can submit verification")
	}
	return s.repo.VerifyEmployer(ctx, userID, inn)
}

func (s *userService) SetPrivacy(ctx context.Context, userID string, isPrivate bool) error {
	return s.repo.UpdatePrivacy(ctx, userID, isPrivate)
}

func (s *userService) SetAvatar(ctx context.Context, userID string, role lib.UserRole, url string) error {
	return s.repo.UpdateAvatar(ctx, userID, role, url)
}
