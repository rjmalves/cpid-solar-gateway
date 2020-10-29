package api

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rjmalves/cpid-solar-gateway/api/models"
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

func checkInverters(t *testing.T, inv1, inv2 models.Inverter) {
	assert.Equal(t, inv1.ID, inv2.ID)
	assert.Equal(t, inv1.Serial, inv2.Serial)
	assert.Equal(t, inv1.Power, inv2.Power)
	assert.Equal(t, inv1.Voltage, inv2.Voltage)
	assert.Equal(t, inv1.Frequency, inv2.Frequency)
	assert.Equal(t, inv1.Communication, inv2.Communication)
	assert.Equal(t, inv1.Status, inv2.Status)
	assert.Equal(t, inv1.Switch, inv2.Switch)
	assert.Equal(t, inv1.EnergyToday, inv2.EnergyToday)
	assert.Equal(t, inv1.EnergyThisMonth, inv2.EnergyThisMonth)
	assert.Equal(t, inv1.EnergyThisYear, inv2.EnergyThisYear)
	assert.Equal(t, inv1.TotalEnergy, inv2.TotalEnergy)
}

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
