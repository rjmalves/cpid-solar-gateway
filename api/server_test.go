package api

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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

func TestMain(m *testing.M) {

	if err := s.Initialize(os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("MODE")); err != nil {
		log.Fatalf("Error initializing the service: %v", err)
	}

	ret := m.Run()

	if err := s.Terminate(); err != nil {
		log.Fatalf("Error terminating the service: %v", err)
	}

	os.Exit(ret)
}
