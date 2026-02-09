package grpc

import (
	"context"
	"errors"
	"strings"

	"buf.build/go/protovalidate"
	userv1 "github.com/auroban/gochat-be/apis/proto/user/v1"
	"github.com/auroban/gochat-be/user-service/internal/domainerror"
	"github.com/auroban/gochat-be/user-service/internal/models"
	"github.com/auroban/gochat-be/user-service/internal/repository"
	"github.com/auroban/gochat-be/user-service/internal/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var logger = log.WithField("package", "grpc")

type UserServer struct {
	userv1.UnimplementedUserServiceServer
	svc       *service.UserService
	validator protovalidate.Validator
}

func NewUserServer(svc *service.UserService) *UserServer {
	validator, _ := protovalidate.New()
	return &UserServer{svc: svc, validator: validator}
}

func (s *UserServer) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	serviceReq := models.ReqCreateUser{
		Username:     req.GetUsername(),
		Password:     req.GetPassword(),
		Email:        req.GetEmail(),
		DisplayName:  req.GetDisplayName(),
		AuthProvider: req.GetAuthProvider(),
		Type:         req.GetType(),
	}

	user, err := s.svc.CreateNewUser(ctx, serviceReq)
	if err != nil {
		if errors.Is(err, domainerror.ErrUserAlreadyExists) {
			logger.WithError(err).Warn("user already exists")
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		logger.WithError(err).Error("failed to create user")
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &userv1.CreateUserResponse{
		Username:     user.Username,
		Email:        user.Email,
		DisplayName:  user.DisplayName,
		AuthProvider: user.AuthProvider,
		Type:         user.Type,
		Status:       mapUserStatus(user.Status),
	}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	var userErr error
	var user *repository.User
	if strings.TrimSpace(req.GetUsername()) != "" {
		user, userErr = s.svc.GetUserByUsername(ctx, req.GetUsername())
	} else {
		user, userErr = s.svc.GetUserByEmail(ctx, req.GetEmail())
	}

	if userErr != nil {
		if errors.Is(userErr, domainerror.ErrNoUserFound) {
			logger.WithError(userErr).Warn("user not found")
			return nil, status.Error(codes.NotFound, userErr.Error())
		}
		logger.WithError(userErr).Error("failed to fetch user")
		return nil, status.Error(codes.Internal, "failed to fetch user")
	}

	return &userv1.GetUserResponse{
		Username:     user.Username,
		Email:        user.Email,
		DisplayName:  user.DisplayName,
		AuthProvider: user.AuthProvider,
		Type:         user.Type,
		Status:       mapUserStatus(user.Status),
	}, nil
}

func (s *UserServer) validateRequest(msg proto.Message) error {
	if s.validator == nil {
		logger.Error("validator not initialized")
		return status.Error(codes.Internal, "validator not initialized")
	}
	if err := s.validator.Validate(msg); err != nil {
		logger.WithError(err).Warn("validation failed")
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}

func mapUserStatus(status string) userv1.UserStatus {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "active":
		return userv1.UserStatus_USER_STATUS_ACTIVE
	case "inactive":
		return userv1.UserStatus_USER_STATUS_INACTIVE
	case "suspended":
		return userv1.UserStatus_USER_STATUS_SUSPENDED
	default:
		return userv1.UserStatus_USER_STATUS_UNSPECIFIED
	}
}
