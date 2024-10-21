package handlers

import (
    "context"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "trail-finder/db"

    "github.com/stretchr/testify/assert"
)

// Setup test database
func setupTestDB(t *testing.T) {
    testDBConnString := "postgresql://postgres:mysecretpassword@localhost:5432/trails_db_test"
    err := db.InitDB(testDBConnString)
    if err != nil {
        t.Fatalf("Failed to initialize test database: %v", err)
    }

    // Clear the trails table before each test
    _, err = db.DbConn.Exec(context.Background(), "TRUNCATE TABLE trails RESTART IDENTITY CASCADE")
    if err != nil {
        t.Fatalf("Failed to clear trails table: %v", err)
    }
}

// Tear down test database
func tearDownTestDB(t *testing.T) {
    db.CloseDB()
}

// Test LoadTrailsFromRequest
func TestLoadTrailsFromRequest(t *testing.T) {
    setupTestDB(t)
    defer tearDownTestDB(t)

    reqBody := `{"file_path": "../BoulderTrailHeads.csv"}`
    req := httptest.NewRequest(http.MethodPost, "/load", strings.NewReader(reqBody))
    w := httptest.NewRecorder()

    LoadTrailsFromRequest(w, req)

    // Assert the response code and body
    assert.Equal(t, http.StatusOK, w.Code, "expected status OK")
    assert.Contains(t, w.Body.String(), "Trails loaded successfully", "expected response to contain 'Trails loaded successfully'")
}

// Test GetTrails
func TestGetTrails(t *testing.T) {
    setupTestDB(t)
    defer tearDownTestDB(t)

    // Insert mock data with non-nullable fields
    _, err := db.DbConn.Exec(context.Background(), `
        INSERT INTO trails 
        (fid, name, restrooms, picnic, fishing, type, difficulty, access_type, th_leash, bike_trail, horse_trail, fee, recycle_bin, grills, bike_rack, dog_tube) 
        VALUES 
        (1, 'Test Trail', 'yes', 'yes', 'no', 'hiking', 'easy', 'TH', 'Yes', 'yes', 'possible', 'no', 'yes', '2', 'yes', '1')
    `)
    if err != nil {
        t.Fatalf("Failed to insert mock data: %v", err)
    }

    // Make a request to the /trails endpoint
    req := httptest.NewRequest(http.MethodGet, "/trails?page=1&limit=5", nil)
    w := httptest.NewRecorder()

    // Call the handler function
    GetTrails(w, req)

    // Assert the response code
    assert.Equal(t, http.StatusOK, w.Code, "expected status OK")

    // Assert the response body contains the mock data
    responseBody := w.Body.String()
    assert.Contains(t, responseBody, "Test Trail", "expected response to contain 'Test Trail'")
}
