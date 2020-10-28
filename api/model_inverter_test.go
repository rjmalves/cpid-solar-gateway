package api

import (
	"context"
	"log"
	"testing"

	"github.com/rjmalves/cpid-solar-gateway/api/models"
	"github.com/rjmalves/cpid-solar-gateway/api/seed"
	"github.com/stretchr/testify/assert"
)

func TestListInverters(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadInverters(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Verifies the inverters in DB
	invs, err := models.ListInverters(s.DB)
	if err != nil {
		t.Errorf("Error while listing inverters in DB: %v\n", err)
		return
	}
	assert.Equal(t, 3, len(invs))
}

func TestReadMissingInverterBySerial(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Tries to read an inverter with the empty collection
	i := models.Inverter{
		Serial: "INVERTER1",
	}
	if err := i.ReadInverter(s.DB); err == nil {
		t.Errorf("Found an inverter that should not exist in DB\n")
		return
	}
}

func TestReadExistingInverterBySerial(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadInverters(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Tries to read an inverter by serial
	i := models.Inverter{
		Serial: "INVERTER1",
	}
	if err := i.ReadInverter(s.DB); err != nil {
		t.Errorf("Error while reading an inverter by serial: %v\n", err)
		return
	}
	assert.Equal(t, "INVERTER1", i.Serial)
	assert.Equal(t, 500.0, i.Power)
}

func TestCreatingRepeatedInverter(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadInverters(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Tries to add an repeated inverter to DB
	i := models.Inverter{
		Serial: "INVERTER1",
	}
	if _, err := i.AddInverterToDB(s.DB); err == nil {
		t.Errorf("Should have failed while adding an repeated inverter\n")
		return
	}
}

func TestCreatingNewInverter(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadInverters(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Tries to create a new inverter
	i := models.Inverter{
		Serial: "INVERTER4",
	}
	if _, err := i.AddInverterToDB(s.DB); err != nil {
		t.Errorf("Failed while adding a new inverter to DB: %v\n", err)
		return
	}
	// List the existing inverters and checks the amount
	invs, _ := models.ListInverters(s.DB)
	assert.Equal(t, 4, len(invs))
}
