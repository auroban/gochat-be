package server

import (
	"fmt"
	"net/http"

	"github.com/auroban/gochat-be/user-service/internal/config"
	"github.com/auroban/gochat-be/user-service/internal/db"
	"github.com/auroban/gochat-be/user-service/internal/handler"
	"github.com/auroban/gochat-be/user-service/internal/repository"
	"github.com/auroban/gochat-be/user-service/internal/router"
	"github.com/auroban/gochat-be/user-service/internal/service"
	"github.com/go-playground/validator/v10"

	log "github.com/sirupsen/logrus"
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
	dbConn := db.Connect(*config)
	db.RunMigration(databaseURL)
	validator := validator.New()
	userRepository := repository.NewUserRepository(dbConn)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(validator, userService)
	router := router.NewRouter(userHandler)
	addr := fmt.Sprintf(":%d", config.Server.Port)
	logger.Infof("Starting server on port: %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}

}
