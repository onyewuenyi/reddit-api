package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

app := &App{}

func TestHealthCheckHandler(t *testing.T) {
	

	// Create http req
	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	// NewRecorder implements ReponseWriter interface to record the HTTP res
	rr := httptest.NewRecorder()

	// Wrap handler fn with http.HandlerFunc to create a http.Handler to pass to ServeHTTP.
	// ServeHTTP expects an input of http.Handler
	handler := http.HandlerFunc(app.healthCheck)
	handler.ServeHTTP(rr, req)

	//
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("healthCheck handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := ""
	if rr.Body.String() != expected {
		t.Errorf("healthCheck handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}


// func TestcreatePostHandler(t *testing.T) {
// 	req, err := http.NewRequest("POST")
// }