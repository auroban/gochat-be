package service

import (
	"context"

	"github.com/auroban/gochat-be/profile-service/internal/models"
	"github.com/auroban/gochat-be/profile-service/internal/repository"
	"github.com/auroban/gochat-be/profile-service/internal/security"
)

type ProfileService struct {
	repository *repository.ProfileRepository
	hasher     security.PasswordHasher
}

func NewProfileService(repository *repository.ProfileRepository, passwordHasher security.PasswordHasher) *ProfileService {
	return &ProfileService{
		repository: repository,
		hasher:     passwordHasher,
	}
}

func (s *ProfileService) CreateNewProfile(ctx context.Context, req models.ReqCreateProfile) (*repository.Profile, error) {
	profile := req.ToProfile()
	hashedPW, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	profile.PasswordHash = hashedPW
	profile.Status = "active"
	return s.repository.CreateProfile(ctx, profile)
}

func (s *ProfileService) GetProfileByID(ctx context.Context, id int) (*repository.Profile, error) {
	return s.repository.FindById(ctx, id)
}

func (s *ProfileService) GetProfileByUsername(ctx context.Context, username string) (*repository.Profile, error) {
	return s.repository.FindByUsername(ctx, username)
}

func (s *ProfileService) GetProfileByEmail(ctx context.Context, email string) (*repository.Profile, error) {
	return s.repository.FindByEmail(ctx, email)
}
