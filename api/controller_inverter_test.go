package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/rjmalves/cpid-solar-gateway/api/models"
	"github.com/rjmalves/cpid-solar-gateway/api/seed"
	"github.com/stretchr/testify/assert"
)

func TestAPIv1ListEmptyInverters(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Makes the request
	res, err := makeRequest(s.Router, "GET", "/api/v1/inverters", nil)
	if err != nil {
		t.Errorf("Error while listing inverters: %v\n", err)
	}
	resBody := []models.Inverter{}
	err = json.Unmarshal(res.Body.Bytes(), &resBody)
	if err != nil {
		t.Errorf("Error converting response to JSON: %v\n", err)
	}
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, 0, len(resBody))
}

func TestAPIv1ListInverters(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadInverters(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Makes the request
	res, err := makeRequest(s.Router, "GET", "/api/v1/inverters", nil)
	if err != nil {
		t.Errorf("Error while listing inverters: %v\n", err)
		return
	}
	data := []models.Inverter{}
	err = json.Unmarshal(res.Body.Bytes(), &data)
	if err != nil {
		t.Errorf("Error converting response to JSON: %v\n", err)
		return
	}
	assert.Equal(t, http.StatusOK, res.Code)
	// Reads the inverters directly from the DB and checks the data
	invs, err := models.ListInverters(s.DB)
	if err != nil {
		t.Errorf("Error while listing inverters in DB: %v\n", err)
		return
	}
	assert.Equal(t, len(invs), len(data))
	for i := range invs {
		assert.Equal(t, *invs[i], data[i])
	}
}

func TestAPIv1GetInverterBySerialMissingSerial(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadInverters(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Makes the request
	res, err := makeRequest(s.Router, "GET", "/api/v1/inverters/0", nil)
	if err != nil {
		t.Errorf("Error while listing inverters: %v\n", err)
		return
	}
	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestAPIv1GetInverterBySerialSuccess(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadInverters(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Makes the request
	res, err := makeRequest(s.Router, "GET", "/api/v1/inverters/INVERTER1", nil)
	if err != nil {
		t.Errorf("Error while listing inverters: %v\n", err)
		return
	}
	data := models.Inverter{}
	err = json.Unmarshal(res.Body.Bytes(), &data)
	if err != nil {
		t.Errorf("Error converting response to JSON: %v\n", err)
		return
	}
	assert.Equal(t, http.StatusOK, res.Code)
	// Reads the inverter directly from the DB and checks the data
	i := models.Inverter{
		Serial: "INVERTER1",
	}
	err = i.ReadInverter(s.DB)
	if err != nil {
		t.Errorf("Error while reading inverter from DB: %v\n", err)
		return
	}
	assert.Equal(t, i, data)
}
