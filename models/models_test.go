package models

import (
    "context"
    "testing"
    "trail-finder/db"

    "github.com/stretchr/testify/assert"
)

func TestCreateTable(t *testing.T) {
    // Initialize the database connection for testing
    testDBConnString := "postgresql://postgres:mysecretpassword@localhost:5432/trails_db_test"
    err := db.InitDB(testDBConnString)
    if err != nil {
        t.Fatalf("Failed to initialize test database: %v", err)
    }
    defer db.CloseDB()

    // Clear the trails table before running the test
    _, err = db.DbConn.Exec(context.Background(), "TRUNCATE TABLE trails RESTART IDENTITY CASCADE")
    if err != nil {
        t.Fatalf("Failed to clear trails table: %v", err)
    }

    // Call the CreateTable function to test
    err = CreateTable(db.DbConn)
    if err != nil {
        t.Fatalf("Failed to create table: %v", err)
    }

    // Verify the table was created successfully
    var exists bool
    err = db.DbConn.QueryRow(context.Background(), `
        SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_name = 'trails'
        )
    `).Scan(&exists)
    if err != nil {
        t.Fatalf("Failed to check if table exists: %v", err)
    }

    assert.True(t, exists, "Expected 'trails' table to exist")
}
