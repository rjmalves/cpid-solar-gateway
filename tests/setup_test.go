package tests

import (
	"log"
	"os"
	"testing"

	"github.com/rjmalves/cpid-solar-gateway/api/controllers"
)

var s = controllers.Server{}

func TestMain(m *testing.M) {

	// Initializes the application
	if err := s.Initialize(os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("MODE")); err != nil {
		log.Fatalf("Error initializing the service: %v", err)
	}

	os.Exit(m.Run())
}
