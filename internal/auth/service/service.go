package service

import (
	"context"
	"errors"
	"springboard/internal/auth/dto"
	"springboard/internal/auth/repository"
	"springboard/internal/lib"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	RefreshTokens(ctx context.Context, refreshToken string) (dto.RefreshTokenResponse, error)
}

type authService struct {
	repo       repository.AuthRepository
	jwtManager *lib.JWTManager
}

func NewAuthService(repo repository.AuthRepository, jwtManager *lib.JWTManager) AuthService {
	return &authService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

// register a new user
func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	passwordHash, err := lib.GenerateHashByPassword(req.Password)
	if err != nil {
		return dto.RegisterResponse{}, err
	}

	userEntity := lib.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         string(req.Role),
		FullName:     req.FullName,
	}

	createdUser, err := s.repo.CreateUser(ctx, userEntity)
	if err != nil {
		return dto.RegisterResponse{}, err
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(
		createdUser.ID,
		createdUser.Email,
		createdUser.Role,
	)
	if err != nil {
		return dto.RegisterResponse{}, err
	}

	return dto.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := s.repo.LoginUser(ctx, req.Email)
	if err != nil {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(
		user.ID,
		user.Email,
		user.Role,
	)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshTokens(ctx context.Context, refreshToken string) (dto.RefreshTokenResponse, error) {
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return dto.RefreshTokenResponse{}, err
	}

	accessToken, newRefreshToken, err := s.jwtManager.GenerateTokenPair(
		claims.UserID,
		claims.Email,
		claims.Role,
	)
	if err != nil {
		return dto.RefreshTokenResponse{}, err
	}

	return dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
