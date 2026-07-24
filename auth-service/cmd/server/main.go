package main

import (
	"github.com/auroban/gochat-be/auth-service/internal/config"
	"github.com/auroban/gochat-be/auth-service/internal/logger"
	"github.com/auroban/gochat-be/auth-service/internal/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := config.Load()
	logger.Init(config.Log.Level)
	log.Debugf("Loaded config: [%v]", config)
	server.Start(config)
}
