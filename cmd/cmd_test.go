package cmd

import (
    "context"
    "os"
    "testing"
    "trail-finder/db"
    "github.com/stretchr/testify/assert"
)

func TestLoadCommand(t *testing.T) {
    // Set up environment variable for the test database
    os.Setenv("DB_CONN_STRING", "postgresql://postgres:mysecretpassword@localhost:5432/trails_db_test")

    // Initialize the database connection
    err := db.InitDB(os.Getenv("DB_CONN_STRING"))
    if err != nil {
        t.Fatalf("Failed to initialize database: %v", err)
    }

    // Begin a new transaction for the test
    tx, err := db.DbConn.Begin(context.Background())
    if err != nil {
        t.Fatalf("Failed to start transaction: %v", err)
    }
    defer tx.Rollback(context.Background()) // Ensure rollback after the test to avoid permanent changes

    // Set up the file flag for the test CSV file
    err = loadCmd.Flags().Set("file", "../BoulderTrailHeads.csv")
    if err != nil {
        t.Fatalf("Failed to set file flag: %v", err)
    }

    // Execute the load command
    loadCSV(loadCmd, []string{})

    // Assertions: Check if data is loaded into the database
    var count int
    err = tx.QueryRow(context.Background(), "SELECT COUNT(*) FROM trails").Scan(&count)
    if err != nil {
        t.Fatalf("Failed to count rows in trails table: %v", err)
    }

    // Verify that some rows are loaded (the row count should be greater than 0)
    assert.Greater(t, count, 0, "Expected some rows in trails table, got %d", count)

    // Clean up the test data within the transaction (rolled back after the test)
    _, err = tx.Exec(context.Background(), "DELETE FROM trails")
    if err != nil {
        t.Fatalf("Failed to clean up trails table: %v", err)
    }

    // Explicitly commit the transaction (only if needed)
    err = tx.Commit(context.Background())
    if err != nil {
        t.Fatalf("Failed to commit transaction: %v", err)
    }

    // Explicitly close the database connection at the end of the test
    db.CloseDB()
}
