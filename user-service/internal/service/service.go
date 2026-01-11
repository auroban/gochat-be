package service

import (
	"context"

	"github.com/auroban/gochat-be/user-service/internal/models"
	"github.com/auroban/gochat-be/user-service/internal/repository"
	"github.com/auroban/gochat-be/user-service/internal/security"
)

type UserService struct {
	repository *repository.UserRepository
	hasher     security.PasswordHasher
}

func NewUserService(repository *repository.UserRepository, passwordHasher security.PasswordHasher) *UserService {
	return &UserService{
		repository: repository,
		hasher:     passwordHasher,
	}
}

func (s *UserService) CreateNewUser(ctx context.Context, req models.ReqCreateUser) (*repository.User, error) {
	user := req.ToUser()
	hashedPW, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = hashedPW
	return s.repository.CreateUser(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*repository.User, error) {
	return s.repository.FindById(ctx, id)
}
