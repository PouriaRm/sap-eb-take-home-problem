package tests

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"
    "trail-finder/handlers"
    "trail-finder/db"

    "github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) {
    // Set up the environment variable for the test database
    testDBConnString := "postgresql://postgres:mysecretpassword@localhost:5432/trails_db_test"
    err := db.InitDB(testDBConnString)
    if err != nil {
        t.Fatalf("Failed to initialize test database: %v", err)
    }

    // Insert mock data into the trails table
     _, err = db.DbConn.Exec(context.Background(), `
        INSERT INTO trails 
        (fid, name, restrooms, picnic, fishing, type, difficulty, access_type, th_leash, bike_trail, horse_trail, fee, recycle_bin, grills, bike_rack, dog_tube) 
        VALUES 
        (1, 'Test Trail', 'yes', 'yes', 'no', 'hiking', 'easy', 'TH', 'Yes', 'yes', 'possible', 'no', 'yes', '2', 'yes', '1')
    `)
    if err != nil {
        t.Fatalf("Failed to insert mock data: %v", err)
    }
}

func tearDownTestDB(t *testing.T) {
    // Clean up the trails table and close the database connection
    _, err := db.DbConn.Exec(context.Background(), "DELETE FROM trails")
    if err != nil {
        t.Fatalf("Failed to clean up trails table: %v", err)
    }
    db.CloseDB()
}

func TestAPIEndpoints(t *testing.T) {
    // Set up the test database
    setupTestDB(t)
    defer tearDownTestDB(t)

    // Create a request to the /trails endpoint
    req := httptest.NewRequest(http.MethodGet, "/trails?page=1&limit=5", nil)
    w := httptest.NewRecorder()

    // Call the GetTrails handler
    handlers.GetTrails(w, req)

    // Check the response status code
    assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

    // Check that the response contains the mock trail data
    responseBody := w.Body.String()
    assert.Contains(t, responseBody, "Test Trail", "Expected response to contain 'Test Trail'")
}
