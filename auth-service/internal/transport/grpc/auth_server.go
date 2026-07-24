package grpc

import (
	"context"

	"buf.build/go/protovalidate"
	authv1 "github.com/auroban/gochat-be/apis/proto/auth/v1"
	"github.com/auroban/gochat-be/auth-service/internal/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var logger = log.WithField("package", "grpc")

type AuthServer struct {
	authv1.UnimplementedAuthServiceServer
	svc       *service.AuthService
	validator protovalidate.Validator
}

func NewAuthServer(svc *service.AuthService) *AuthServer {
	validator, _ := protovalidate.New()
	return &AuthServer{svc: svc, validator: validator}
}

// HealthCheck is a placeholder RPC so the gRPC server builds and serves.
// Real auth RPCs (Register, Login, GetUser, ...) replace it later.
func (s *AuthServer) HealthCheck(ctx context.Context, req *authv1.HealthCheckRequest) (*authv1.HealthCheckResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}
	return &authv1.HealthCheckResponse{Status: "ok"}, nil
}

func (s *AuthServer) validateRequest(msg proto.Message) error {
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
