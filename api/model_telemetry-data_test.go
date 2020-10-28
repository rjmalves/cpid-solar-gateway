package api

import (
	"context"
	"log"
	"testing"

	"github.com/rjmalves/cpid-solar-gateway/api/models"
	"github.com/rjmalves/cpid-solar-gateway/api/seed"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestListTelemetryData(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshTelemetryDataCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadTelemetryData(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Verifies the data in DB
	tel, err := models.ListTelemetryData(s.DB, bson.M{})
	if err != nil {
		t.Errorf("Error while listing telemetry data in DB: %v\n", err)
		return
	}
	assert.Equal(t, 300, len(tel))
}

func TestListMissingTelemetryDataBySerial(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshTelemetryDataCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Tries to read data with the empty collection
	filter := bson.M{
		"serial": "INVERTER1",
	}
	data, _ := models.ListTelemetryData(s.DB, filter)
	assert.Equal(t, 0, len(data))
}

func TestListExistingTelemetryDataBySerial(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshTelemetryDataCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadTelemetryData(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Tries to read data filtered by serial
	filter := bson.M{
		"serial": "INVERTER1",
	}
	data, err := models.ListTelemetryData(s.DB, filter)
	if err != nil {
		t.Errorf("Error while listing data by serial: %v\n", err)
		return
	}
	assert.NotEqual(t, 0, len(data))
}

func TestCreatingRepeatedTelemetryData(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshTelemetryDataCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadTelemetryData(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Tries to add an repeated inverter to DB
	tel := models.TelemetryData{
		Serial:            "INVERTER1",
		LastTelemetryTime: int64(0),
	}
	if _, err := tel.AddDataToDB(s.DB); err == nil {
		t.Errorf("Should have failed while adding an repeated telemetry data\n")
		return
	}
}

func TestCreatingNewTelemetryData(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshTelemetryDataCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadTelemetryData(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Tries to create a new telemetry entry
	tel := models.TelemetryData{
		Serial:            "INVERTER4",
		LastTelemetryTime: int64(0),
	}
	if _, err := tel.AddDataToDB(s.DB); err != nil {
		t.Errorf("Failed while adding a new telemetry entry to DB: %v\n", err)
		return
	}
	// List the existing telemetry data and checks the amount
	filter := bson.M{
		"serial": "INVERTER4",
	}
	data, _ := models.ListTelemetryData(s.DB, filter)
	assert.Equal(t, 1, len(data))
}
