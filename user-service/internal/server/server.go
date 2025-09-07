package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/auroban/gochat-be/user-service/internal/config"
)

func Start(config *config.Config) {

	log.Printf("Starting server on port: [%d]", config.Server.Port)
	addr := fmt.Sprintf(":%d", config.Server.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
