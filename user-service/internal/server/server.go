package server

import (
	"fmt"
	"net"

	userv1 "github.com/auroban/gochat-be/apis/proto/user/v1"
	"github.com/auroban/gochat-be/user-service/internal/config"
	"github.com/auroban/gochat-be/user-service/internal/db"
	"github.com/auroban/gochat-be/user-service/internal/repository"
	"github.com/auroban/gochat-be/user-service/internal/security"
	"github.com/auroban/gochat-be/user-service/internal/service"
	transportGrpc "github.com/auroban/gochat-be/user-service/internal/transport/grpc"
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
	userService := service.NewUserService(userRepository, security.BcryptHasher{})

	addr := fmt.Sprintf(":%d", config.Server.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalf("Failed to listen on %s: %v", addr, err)
	}

	grpcServer := grpc.NewServer()
	userServer := transportGrpc.NewUserServer(userService)
	userv1.RegisterUserServiceServer(grpcServer, userServer)

	logger.Infof("Starting gRPC server on %s", addr)
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatalf("Failed to start gRPC server: %v", err)
	}
}
