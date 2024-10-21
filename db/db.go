package db

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v4"
)

var DbConn *pgx.Conn

// InitDB initializes the connection to PostgreSQL
func InitDB(connString string) error {
    var err error
    DbConn, err = pgx.Connect(context.Background(), connString)
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    // Create the trails table if it doesn't exist
    if err := CreateTable(); err != nil {
        return fmt.Errorf("failed to create trails table: %w", err)
    }

    fmt.Println("Database initialized successfully.")
    return nil
}

// CreateTable creates the trails table if it doesn't exist
func CreateTable() error {
    createTableQuery := `
        CREATE TABLE IF NOT EXISTS trails (
            fid SERIAL PRIMARY KEY,
            name TEXT,
            restrooms TEXT,
            picnic TEXT,
            fishing TEXT,
			type TEXT,
            bike_trail TEXT,
            horse_trail TEXT,
            difficulty TEXT,
            fee TEXT,
            recycle_bin TEXT,
            grills TEXT,
            bike_rack TEXT,
            dog_tube TEXT
        )
    `

    _, err := DbConn.Exec(context.Background(), createTableQuery)
    if err != nil {
        return fmt.Errorf("failed to create trails table: %w", err)
    }

    fmt.Println("Trails table created or already exists.")
    return nil
}

// CloseDB closes the PostgreSQL connection
func CloseDB() {
    if DbConn != nil {
        err := DbConn.Close(context.Background())
        if err != nil {
            fmt.Println("Error closing database connection:", err)
        } else {
            fmt.Println("Database connection closed.")
        }
    }
}

// IsTableEmpty checks if the trails table is empty
func IsTableEmpty() (bool, error) {
    var count int
    err := DbConn.QueryRow(context.Background(), "SELECT COUNT(*) FROM trails").Scan(&count)
    if err != nil {
        return false, fmt.Errorf("failed to count trails: %w", err)
    }
    return count == 0, nil
}
