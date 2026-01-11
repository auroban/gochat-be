package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/auroban/gochat-be/user-service/internal/models"
	"github.com/auroban/gochat-be/user-service/internal/repository"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository *repository.UserRepository
}

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrNoUserFound = errors.New("no user found")

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) CreateNewUser(ctx context.Context, req models.ReqCreateUser) (*repository.User, error) {
	user := req.ToUser()
	hashedPW, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = hashedPW
	createdUser, err := s.repository.CreateUser(ctx, user)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return nil, fmt.Errorf("%w: %s", ErrUserAlreadyExists, req.Username)
			}
		}
		return nil, err
	}
	return createdUser, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*repository.User, error) {
	user, err := s.repository.FindById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: %d", ErrNoUserFound, id)
		}
		return nil, err
	}
	return user, nil
}

func hashPassword(rawPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
