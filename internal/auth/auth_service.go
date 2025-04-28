package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinxinyu/go_backend/internal/api"
	"github.com/jinxinyu/go_backend/internal/models"
	"github.com/jinxinyu/go_backend/internal/storage"
	"github.com/jinxinyu/go_backend/internal/utils"
)

type Service struct {
	userRepo     storage.UserRepository
	timeout      time.Duration
	tokenmaker   utils.ToKenGenerator
	hashPassword utils.HashedPassword
	//To be soon added: emailservice
}

func NewService(userRepo storage.UserRepository, tokenmaker utils.ToKenGenerator, hashPassword utils.HashedPassword) *Service {
	return &Service{
		userRepo:     userRepo,
		tokenmaker:   tokenmaker,
		hashPassword: hashPassword,
	}
}

func (s *Service) RegisterUser(ctx context.Context, req *api.RegisterRequest) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	hashedPassword, err := s.hashPassword.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) LoginUser(ctx context.Context, req *api.LoginRequest) (string, *models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	match, err := s.hashPassword.Compare(req.Password, user.Password)
	if err != nil {
		return "", nil, fmt.Errorf("failed to compare password: %w", err)
	}

	if !match {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	token, err := s.tokenmaker.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}
