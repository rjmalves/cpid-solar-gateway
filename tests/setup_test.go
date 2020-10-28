package tests

import (
	"log"
	"os"
	"testing"

	"github.com/rjmalves/cpid-solar-gateway/api/controllers"
	"github.com/rjmalves/cpid-solar-gateway/api/seed"
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

	// Seeds the DB for testing data
	if err := seed.LoadInverters(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	if err := seed.LoadTelemetryData(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}

	os.Exit(m.Run())
}
