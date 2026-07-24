package grpc

import (
	"context"
	"errors"
	"strings"

	"buf.build/go/protovalidate"
	profilev1 "github.com/auroban/gochat-be/apis/proto/profile/v1"
	"github.com/auroban/gochat-be/profile-service/internal/domainerror"
	"github.com/auroban/gochat-be/profile-service/internal/models"
	"github.com/auroban/gochat-be/profile-service/internal/repository"
	"github.com/auroban/gochat-be/profile-service/internal/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var logger = log.WithField("package", "grpc")

type ProfileServer struct {
	profilev1.UnimplementedProfileServiceServer
	svc       *service.ProfileService
	validator protovalidate.Validator
}

func NewProfileServer(svc *service.ProfileService) *ProfileServer {
	validator, _ := protovalidate.New()
	return &ProfileServer{svc: svc, validator: validator}
}

func (s *ProfileServer) CreateProfile(ctx context.Context, req *profilev1.CreateProfileRequest) (*profilev1.CreateProfileResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	serviceReq := models.ReqCreateProfile{
		Username:     req.GetUsername(),
		Password:     req.GetPassword(),
		Email:        req.GetEmail(),
		DisplayName:  req.GetDisplayName(),
		AuthProvider: req.GetAuthProvider(),
		Type:         req.GetType(),
	}

	profile, err := s.svc.CreateNewProfile(ctx, serviceReq)
	if err != nil {
		if errors.Is(err, domainerror.ErrProfileAlreadyExists) {
			logger.WithError(err).Warn("profile already exists")
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		logger.WithError(err).Error("failed to create profile")
		return nil, status.Error(codes.Internal, "failed to create profile")
	}

	return &profilev1.CreateProfileResponse{
		Username:     profile.Username,
		Email:        profile.Email,
		DisplayName:  profile.DisplayName,
		AuthProvider: profile.AuthProvider,
		Type:         profile.Type,
		Status:       mapProfileStatus(profile.Status),
	}, nil
}

func (s *ProfileServer) GetProfile(ctx context.Context, req *profilev1.GetProfileRequest) (*profilev1.GetProfileResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	var profileErr error
	var profile *repository.Profile
	if strings.TrimSpace(req.GetUsername()) != "" {
		profile, profileErr = s.svc.GetProfileByUsername(ctx, req.GetUsername())
	} else {
		profile, profileErr = s.svc.GetProfileByEmail(ctx, req.GetEmail())
	}

	if profileErr != nil {
		if errors.Is(profileErr, domainerror.ErrNoProfileFound) {
			logger.WithError(profileErr).Warn("profile not found")
			return nil, status.Error(codes.NotFound, profileErr.Error())
		}
		logger.WithError(profileErr).Error("failed to fetch profile")
		return nil, status.Error(codes.Internal, "failed to fetch profile")
	}

	return &profilev1.GetProfileResponse{
		Username:     profile.Username,
		Email:        profile.Email,
		DisplayName:  profile.DisplayName,
		AuthProvider: profile.AuthProvider,
		Type:         profile.Type,
		Status:       mapProfileStatus(profile.Status),
	}, nil
}

func (s *ProfileServer) validateRequest(msg proto.Message) error {
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

func mapProfileStatus(status string) profilev1.ProfileStatus {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "active":
		return profilev1.ProfileStatus_PROFILE_STATUS_ACTIVE
	case "inactive":
		return profilev1.ProfileStatus_PROFILE_STATUS_INACTIVE
	case "suspended":
		return profilev1.ProfileStatus_PROFILE_STATUS_SUSPENDED
	default:
		return profilev1.ProfileStatus_PROFILE_STATUS_UNSPECIFIED
	}
}
