package main

import (
	"github.com/auroban/gochat-be/user-service/internal/config"
	"github.com/auroban/gochat-be/user-service/internal/logger"
	"github.com/auroban/gochat-be/user-service/internal/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := config.Load()
	logger.Init(config.Log.Level)
	log.Debugf("Loaded config: [%v]", config)
	server.Start(config)
}
