package api

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rjmalves/cpid-solar-gateway/api/models"
	"github.com/rjmalves/cpid-solar-gateway/api/seed"
	"github.com/stretchr/testify/assert"
)

func makeRequest(r http.Handler, method, path string, body io.Reader) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr, nil
}

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
		assert.Equal(t, invs[i].ID, data[i].ID)
		assert.Equal(t, invs[i].Serial, data[i].Serial)
		assert.Equal(t, invs[i].Power, data[i].Power)
		assert.Equal(t, invs[i].Voltage, data[i].Voltage)
		assert.Equal(t, invs[i].Frequency, data[i].Frequency)
		assert.Equal(t, invs[i].Communication, data[i].Communication)
		assert.Equal(t, invs[i].Status, data[i].Status)
		assert.Equal(t, invs[i].Switch, data[i].Switch)
		assert.Equal(t, invs[i].EnergyToday, data[i].EnergyToday)
		assert.Equal(t, invs[i].EnergyThisMonth, data[i].EnergyThisMonth)
		assert.Equal(t, invs[i].EnergyThisYear, data[i].EnergyThisYear)
		assert.Equal(t, invs[i].TotalEnergy, data[i].TotalEnergy)
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
	assert.Equal(t, i.ID, data.ID)
	assert.Equal(t, i.Serial, data.Serial)
	assert.Equal(t, i.Power, data.Power)
	assert.Equal(t, i.Voltage, data.Voltage)
	assert.Equal(t, i.Frequency, data.Frequency)
	assert.Equal(t, i.Communication, data.Communication)
	assert.Equal(t, i.Status, data.Status)
	assert.Equal(t, i.Switch, data.Switch)
	assert.Equal(t, i.EnergyToday, data.EnergyToday)
	assert.Equal(t, i.EnergyThisMonth, data.EnergyThisMonth)
	assert.Equal(t, i.EnergyThisYear, data.EnergyThisYear)
	assert.Equal(t, i.TotalEnergy, data.TotalEnergy)
}
