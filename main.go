package main

import (
    "flag"
    "net/http"
    "os"
    "os/signal"
    "os/exec"
    "syscall"
    "trail-finder/handlers"
    "github.com/sirupsen/logrus"
    "github.com/joho/godotenv"
    "trail-finder/db"
)

func main() {
    // Initialize logger
    handlers.InitLogger()

    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        logrus.Warn("Warning: .env file not found or failed to load.")
    }

    // Define command-line flags
    startServer := flag.Bool("server", false, "Start the trail server")
    dbConnString := flag.String("db", os.Getenv("DB_CONN_STRING"), "Postgres connection string")
    runMigrations := flag.Bool("migrate", false, "Run database migrations")
    flag.Parse()

    if *runMigrations {
        if err := db.InitDB(*dbConnString); err != nil {
            logrus.Fatalf("Failed to initialize database: %v", err)
        }
        defer db.CloseDB()

        // Run the migrations
        cmd := exec.Command("trail-cli", "migrate")
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        err := cmd.Run()
        if err != nil {
            logrus.Fatalf("Failed to run migrations: %v", err)
        }
        logrus.Info("Migrations completed.")
        return
    }


    if *startServer {
        // Initialize the database connection
        if err := db.InitDB(*dbConnString); err != nil {
            logrus.Fatalf("Failed to initialize database: %v", err)
        }
        // Ensure the database is closed on exit
        defer db.CloseDB()

        // Load default data if the table is empty
        if err := handlers.LoadDefaultData("./BoulderTrailHeads.csv"); err != nil {
            logrus.Warnf("Failed to load default data: %v", err)
        }

        // Register the /load endpoint
        http.HandleFunc("/load", handlers.LoadTrailsFromRequest)

        // Register the /trails endpoint
        http.HandleFunc("/trails", handlers.GetTrails)

        // Handle graceful shutdown
        go func() {
            c := make(chan os.Signal, 1)
            signal.Notify(c, os.Interrupt, syscall.SIGTERM)
            <-c
            logrus.Info("Received termination signal, shutting down...")
            os.Exit(0)
        }()

        // Start the server
        logrus.Info("Server running on http://localhost:8080")
        if err := http.ListenAndServe(":8080", nil); err != nil {
            logrus.Fatalf("Server error: %v", err)
        }
    } else {
        logrus.Info("This program is intended to start the server with the --server flag.")
        logrus.Info("Usage: go run main.go --server")
    }
}
