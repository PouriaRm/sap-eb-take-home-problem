package db

import (
    "context"
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) {
    // Set up environment variable for the test database
    os.Setenv("DB_CONN_STRING", "postgresql://postgres:mysecretpassword@localhost:5432/trails_db_test")

    // Initialize the database connection
    err := InitDB(os.Getenv("DB_CONN_STRING"))
    if err != nil {
        t.Fatalf("Failed to initialize test database: %v", err)
    }

    // Clear the trails table before each test to prevent duplicates
    _, err = DbConn.Exec(context.Background(), "TRUNCATE TABLE trails RESTART IDENTITY CASCADE")
    if err != nil {
        t.Fatalf("Failed to clear trails table: %v", err)
    }
}

func tearDownTestDB(t *testing.T) {
    // Clean up the database after tests
    if DbConn != nil {
        err := DbConn.Close(context.Background())
        if err != nil {
            t.Fatalf("Failed to close test database connection: %v", err)
        }
    }
}

func TestDatabaseConnection(t *testing.T) {
    setupTestDB(t)
    defer tearDownTestDB(t)

    // Check if the connection is established
    assert.NotNil(t, DbConn, "Expected DbConn to be initialized")
}

func TestCreateTable(t *testing.T) {
    setupTestDB(t)
    defer tearDownTestDB(t)

    // Verify that the table is created
    var tableName string
    err := DbConn.QueryRow(context.Background(), "SELECT table_name FROM information_schema.tables WHERE table_name='trails'").Scan(&tableName)
    assert.Nil(t, err, "Expected table 'trails' to exist")
    assert.Equal(t, "trails", tableName, "Expected table name to be 'trails'")
}

func TestInsertTrail(t *testing.T) {
    setupTestDB(t)
    defer tearDownTestDB(t)

    // Insert a sample trail into the database
    _, err := DbConn.Exec(context.Background(), `
    INSERT INTO trails 
    (fid, name, restrooms, picnic, fishing, type, difficulty, access_type, th_leash, bike_trail, horse_trail, fee, recycle_bin, grills, bike_rack, dog_tube) 
    VALUES 
    (1, 'Test Trail', 'yes', 'yes', 'no', 'hiking', 'easy', 'TH', 'Yes', 'yes', 'possible', 'no', 'yes', '2', 'yes', '1')
`)
if err != nil {
    t.Fatalf("Failed to insert mock data: %v", err)
}
    assert.Nil(t, err, "Failed to insert trail into the database")

    // Verify that the trail is inserted
    var count int
    err = DbConn.QueryRow(context.Background(), "SELECT COUNT(*) FROM trails WHERE name='Test Trail'").Scan(&count)
    assert.Nil(t, err, "Failed to count trails in the database")
    assert.Equal(t, 1, count, "Expected 1 row in the trails table")
}

func TestCleanUpTrails(t *testing.T) {
    setupTestDB(t)
    defer tearDownTestDB(t)

    // Clean up the trails table after testing
    _, err := DbConn.Exec(context.Background(), "DELETE FROM trails")
    assert.Nil(t, err, "Failed to clean up the trails table")
}
