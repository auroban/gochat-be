package server

import (
	"fmt"
	"net"

	authv1 "github.com/auroban/gochat-be/apis/proto/auth/v1"
	"github.com/auroban/gochat-be/auth-service/internal/config"
	"github.com/auroban/gochat-be/auth-service/internal/db"
	"github.com/auroban/gochat-be/auth-service/internal/repository"
	"github.com/auroban/gochat-be/auth-service/internal/security"
	"github.com/auroban/gochat-be/auth-service/internal/service"
	transportGrpc "github.com/auroban/gochat-be/auth-service/internal/transport/grpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var logger = log.WithField("package", "server")

func Start(config *config.Config) {

	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&search_path=%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName,
		config.Database.SSLMode,
		config.Database.Schema,
	)
	conn := db.Connect(*config)
	db.RunMigration(databaseURL)

	userRepository := repository.NewUserRepository(conn)
	authService := service.NewAuthService(userRepository, security.BcryptHasher{})

	addr := fmt.Sprintf(":%d", config.Server.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalf("Failed to listen on %s: %v", addr, err)
	}

	grpcServer := grpc.NewServer()
	authServer := transportGrpc.NewAuthServer(authService)
	authv1.RegisterAuthServiceServer(grpcServer, authServer)

	logger.Infof("Starting gRPC server on %s", addr)
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatalf("Failed to start gRPC server: %v", err)
	}
}
