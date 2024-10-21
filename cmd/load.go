package cmd

import (
    "os"
    "trail-finder/handlers"
	"trail-finder/db"
    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
    Use:   "load",
    Short: "Load CSV data into the database",
    Long:  `Load the specified CSV file into the database, replacing existing data.`,
    Run:   loadCSV,
}

func init() {
    loadCmd.Flags().StringP("file", "f", "", "Path to the CSV file")
    rootCmd.AddCommand(loadCmd)
}

func loadCSV(cmd *cobra.Command, args []string) {
    // Initialize the logger
    handlers.InitLogger()

    // Load environment variables from the .env file
    if err := godotenv.Load(); err != nil {
        logrus.Warn("Error loading .env file")
    }

    // Get the database connection string from environment variables
    dbConnString := os.Getenv("DB_CONN_STRING")
    if dbConnString == "" {
        logrus.Error("DB_CONN_STRING environment variable is not set")
        return
    }

    // Initialize the database connection
    if err := db.InitDB(dbConnString); err != nil {
        logrus.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.CloseDB() // Ensure the database connection is closed after the operation

    // Get the file flag
    file, _ := cmd.Flags().GetString("file")
    if file == "" {
        logrus.Error("CSV file path must be provided using the --file flag")
        return
    }

    // Check if the file exists
    if _, err := os.Stat(file); os.IsNotExist(err) {
        logrus.Errorf("Specified file does not exist: %s", file)
        return
    }

    // Load the CSV file into the database
    err := handlers.LoadTrails(file)
    if err != nil {
        logrus.Errorf("Error loading CSV file: %v", err)
        return
    }

    logrus.Infof("CSV data loaded successfully from: %s", file)
}
