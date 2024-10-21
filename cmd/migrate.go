package cmd

import (
    "database/sql"
    "os"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"

    _ "github.com/jackc/pgx/v4/stdlib"
)

var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Run database migrations",
    Long:  "Run database migrations to set up or update the database schema",
    Run:   runMigrations,
}

func init() {
    rootCmd.AddCommand(migrateCmd)
}

func runMigrations(cmd *cobra.Command, args []string) {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        logrus.Warn("Warning: .env file not found or failed to load.")
    }

    // Load the database connection string from environment variables
    dbConnString := os.Getenv("DB_CONN_STRING")
    if dbConnString == "" {
        logrus.Error("DB_CONN_STRING environment variable is not set")
        return
    }

    // Create a new *sql.DB connection using the pgx driver
    db, err := sql.Open("pgx", dbConnString)
    if err != nil {
        logrus.Fatalf("Failed to create database connection: %v", err)
    }
    defer db.Close()

    // Verify the connection is valid
    if err := db.Ping(); err != nil {
        logrus.Fatalf("Failed to connect to the database: %v", err)
    }

    // Use the created *sql.DB connection for the migrate library
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        logrus.Fatalf("Failed to create database driver: %v", err)
    }

    // Initialize the migration
    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations", // Path to the migrations folder
        "postgres",          // Database name
        driver,
    )
    if err != nil {
        logrus.Fatalf("Failed to initialize migrations: %v", err)
    }

    // Run the migrations
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        logrus.Fatalf("Failed to apply migrations: %v", err)
    }

    logrus.Info("Migrations applied successfully!")
}
