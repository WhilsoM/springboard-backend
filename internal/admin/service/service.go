package service

import (
	"context"
	"errors"
	"springboard/internal/admin/dto"
	"springboard/internal/admin/repository"

	"golang.org/x/crypto/bcrypt"
)

type AdminService interface {
	CreateCurator(ctx context.Context, data dto.CreateCuratorRequest) error
	CreateTag(ctx context.Context, name string) (dto.Tag, error)
	GetAllTags(ctx context.Context) ([]dto.Tag, error)
	SubmitVerification(ctx context.Context, employerID, inn, companyName string) error
	ModerateVerification(ctx context.Context, requestID, employerID, status string) error
	ForceDeleteOpportunity(ctx context.Context, oppID string) error
}

type adminService struct {
	repo repository.AdminRepository
}

func NewAdminService(repo repository.AdminRepository) AdminService {
	return &adminService{repo: repo}
}

func (s *adminService) CreateCurator(ctx context.Context, data dto.CreateCuratorRequest) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.CreateCurator(ctx, data.Email, string(hash), data.DisplayName)
}

func (s *adminService) CreateTag(ctx context.Context, name string) (dto.Tag, error) {
	if name == "" {
		return dto.Tag{}, errors.New("tag name cannot be empty")
	}
	return s.repo.CreateTag(ctx, name)
}

func (s *adminService) GetAllTags(ctx context.Context) ([]dto.Tag, error) {
	return s.repo.GetAllTags(ctx)
}

func (s *adminService) SubmitVerification(ctx context.Context, employerID, inn, companyName string) error {
	if inn == "" || companyName == "" {
		return errors.New("inn and company_name are required")
	}

	status := "pending"
	if len(inn) == 10 || len(inn) == 12 {
		status = "approved"
	}

	err := s.repo.CreateVerification(ctx, employerID, inn, companyName, status)
	if err != nil {
		return err
	}

	if status == "approved" {
		return s.repo.ApproveEmployer(ctx, employerID)
	}

	return nil
}

func (s *adminService) ModerateVerification(ctx context.Context, requestID, employerID, status string) error {
	if status != "approved" && status != "rejected" {
		return errors.New("invalid status")
	}

	err := s.repo.UpdateVerificationStatus(ctx, requestID, status)
	if err != nil {
		return err
	}

	if status == "approved" {
		return s.repo.ApproveEmployer(ctx, employerID)
	}
	return nil
}

func (s *adminService) ForceDeleteOpportunity(ctx context.Context, oppID string) error {
	return s.repo.DeleteOpportunity(ctx, oppID)
}
