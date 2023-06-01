package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestCreateUser(t *testing.T) {
	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dbx := sqlx.NewDb(db, "sqlmock")

	// Expectations for the database
	mock.ExpectExec("INSERT INTO users").
		WithArgs("Test User").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create request and response recorder
	req, err := http.NewRequest("POST", "/user", strings.NewReader(`{"name":"Test User"}`))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Create an App instance with the mock DB
	app := &App{DB: dbx}

	// Call the handler function
	app.createUser(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Any other tests you want to do...
}
