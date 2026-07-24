package service

import (
	"context"

	"github.com/auroban/gochat-be/auth-service/internal/repository"
	"github.com/auroban/gochat-be/auth-service/internal/security"
)

type AuthService struct {
	repository *repository.UserRepository
	hasher     security.PasswordHasher
}

func NewAuthService(repository *repository.UserRepository, passwordHasher security.PasswordHasher) *AuthService {
	return &AuthService{
		repository: repository,
		hasher:     passwordHasher,
	}
}

func (s *AuthService) GetUserByID(ctx context.Context, id int) (*repository.User, error) {
	return s.repository.FindById(ctx, id)
}

func (s *AuthService) GetUserByUsername(ctx context.Context, username string) (*repository.User, error) {
	return s.repository.FindByUsername(ctx, username)
}

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	return s.repository.FindByEmail(ctx, email)
}
