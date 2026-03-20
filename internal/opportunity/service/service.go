package service

import (
	"context"
	"errors"
	"log"
	"springboard/internal/lib"
	"springboard/internal/opportunity/dto"
	"springboard/internal/opportunity/repository"
	userRepo "springboard/internal/user/repository"
)

type OpportunityService interface {
	Create(ctx context.Context, userID string, data dto.CreateOpportunityRequest) (dto.OpportunityResponse, error)
	GetAll(ctx context.Context, f dto.SearchFilters) ([]dto.OpportunityResponse, error)
	GetOne(ctx context.Context, id string) (dto.OpportunityResponse, error)
	Update(ctx context.Context, id string, userID string, data dto.CreateOpportunityRequest) error
	Delete(ctx context.Context, id string, userID string) error
	GetEmployerOwn(ctx context.Context, employerID string) ([]dto.OpportunityResponse, error)
}

type opportunityService struct {
	repo     repository.OpportunityRepository
	userRepo userRepo.UserRepository
}

func NewOpportunityService(repo repository.OpportunityRepository, ur userRepo.UserRepository) OpportunityService {
	return &opportunityService{repo: repo, userRepo: ur}
}

func (s *opportunityService) Create(ctx context.Context, userID string, data dto.CreateOpportunityRequest) (dto.OpportunityResponse, error) {
	// check if employer is verified
	profile, err := s.userRepo.GetFullUserByID(ctx, userID, lib.RoleEmployer)
	if err != nil {
		return dto.OpportunityResponse{}, err
	}

	p := profile.(lib.EmployerUser)
	if !p.IsVerified {
		return dto.OpportunityResponse{}, errors.New("employer not verified")
	}

	return s.repo.Create(ctx, userID, p.CompanyName, data)
}

func (s *opportunityService) GetAll(ctx context.Context, f dto.SearchFilters) ([]dto.OpportunityResponse, error) {
	return s.repo.Search(ctx, f)
}

func (s *opportunityService) GetOne(ctx context.Context, id string) (dto.OpportunityResponse, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *opportunityService) Update(ctx context.Context, id string, userID string, data dto.CreateOpportunityRequest) error {
	opp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	log.Printf("employer id = %s, userId = %s\n", opp.EmployerID, userID)
	if opp.EmployerID != userID {
		return errors.New("not an owner")
	}
	return s.repo.Update(ctx, id, data)
}

func (s *opportunityService) Delete(ctx context.Context, id string, userID string) error {
	opp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if opp.EmployerID != userID {
		return errors.New("not an owner")
	}
	return s.repo.Delete(ctx, id)
}

func (s *opportunityService) GetEmployerOwn(ctx context.Context, employerID string) ([]dto.OpportunityResponse, error) {
	return s.repo.GetByEmployerID(ctx, employerID)
}
