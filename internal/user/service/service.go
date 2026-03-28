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
	SearchApplicants(ctx context.Context, query string, limit, offset int) ([]lib.ApplicantUser, error)
	SendRequest(ctx context.Context, senderID, receiverID string) error
	HandleContactRequest(ctx context.Context, userID, requestID, status string) error
	GetMyContacts(ctx context.Context, userID string) ([]lib.User, error)
	GetUserProfile(ctx context.Context, targetID string) (any, error)
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

func (s *userService) SearchApplicants(ctx context.Context, query string, limit, offset int) ([]lib.ApplicantUser, error) {
	log.Print("SearchApplicants service start")
	return s.repo.GetApplicants(ctx, query, limit, offset)
}

func (s *userService) SendRequest(ctx context.Context, senderID, receiverID string) error {
	log.Print("SendRequest service start")
	if senderID == receiverID {
		return errors.New("you cannot add yourself to contacts")
	}
	return s.repo.CreateContactRequest(ctx, senderID, receiverID)
}

func (s *userService) HandleContactRequest(ctx context.Context, userID, requestID, status string) error {
	log.Print("HandleContactRequest service start")
	if status != "accepted" && status != "rejected" {
		return errors.New("invalid status: must be accepted or rejected")
	}
	// TODO: add validation receiver id is owner to user id
	return s.repo.UpdateContactStatus(ctx, requestID, status)
}

func (s *userService) GetMyContacts(ctx context.Context, userID string) ([]lib.User, error) {
	log.Print("GetMyContacts service start")
	return s.repo.GetContacts(ctx, userID)
}

func (s *userService) GetUserProfile(ctx context.Context, targetID string) (any, error) {
	return s.repo.GetPublicProfile(ctx, targetID)
}
