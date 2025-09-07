package main

import (
	"github.com/auroban/gochat-be/user-service/internal/config"
	"github.com/auroban/gochat-be/user-service/internal/server"
)

func main() {
	config := config.Load()
	server.Start(config)
}
