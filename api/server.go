package api

import (
	"log"
	"os"

	"github.com/rjmalves/cpid-solar-gateway/api/controllers"
)

var s = controllers.Server{}

// Run : launches the service
func Run() {
	if err := s.Initialize(os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("MODE")); err != nil {
		log.Fatalf("Error initializing the service: %v", err)
	}
	// seed.Load(s.DB)
	if err := s.Run(os.Getenv("SERVICE_PORT")); err != nil {
		log.Fatalf("Error running the service: %v", err)
	}
}
