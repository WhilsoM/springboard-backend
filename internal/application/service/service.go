package service

import (
	"context"
	"errors"
	"springboard/internal/application/dto"
	"springboard/internal/application/repository"
)

type ApplicationService interface {
	Apply(ctx context.Context, userID, role, oppID string, req dto.ApplyRequest) (dto.ApplicationResponse, error)
	GetForEmployer(ctx context.Context, userID, role, oppID string, f dto.PaginationFilters) ([]dto.ApplicationResponse, error)
	UpdateStatus(ctx context.Context, userID, role, appID string, req dto.UpdateStatusRequest) error
}

type applicationService struct {
	repo repository.ApplicationRepository
}

func NewApplicationService(repo repository.ApplicationRepository) ApplicationService {
	return &applicationService{repo: repo}
}

func (s *applicationService) Apply(ctx context.Context, userID, role, oppID string, req dto.ApplyRequest) (dto.ApplicationResponse, error) {
	if role != "applicant" {
		return dto.ApplicationResponse{}, errors.New("only applicants can apply")
	}
	return s.repo.Create(ctx, userID, oppID, req)
}

func (s *applicationService) GetForEmployer(ctx context.Context, userID, role, oppID string, f dto.PaginationFilters) ([]dto.ApplicationResponse, error) {
	if role != "employer" && role != "admin" {
		return nil, errors.New("forbidden")
	}
	return s.repo.GetByOpportunity(ctx, oppID, f.Limit, f.Offset)
}

func (s *applicationService) UpdateStatus(ctx context.Context, userID, role, appID string, req dto.UpdateStatusRequest) error {
	if role != "employer" && role != "admin" {
		return errors.New("forbidden")
	}

	if role == "employer" {
		ok, _ := s.repo.CheckOwnership(ctx, appID, userID)
		if !ok {
			return errors.New("not your opportunity")
		}
	}
	return s.repo.UpdateStatus(ctx, appID, req.Status)
}
