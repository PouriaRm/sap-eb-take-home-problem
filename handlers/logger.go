package handlers

import (
    "github.com/sirupsen/logrus"
    "os"
)

// Logger is the global logger for the application
var Logger = logrus.New()

// InitLogger initializes the logger with custom settings
func InitLogger() {
    // Set log format to JSON
    Logger.SetFormatter(&logrus.JSONFormatter{})

    // Set log output to stdout
    Logger.SetOutput(os.Stdout)

    // Set log level (can be set via environment variable as well)
    logLevel, exists := os.LookupEnv("LOG_LEVEL")
    if exists {
        level, err := logrus.ParseLevel(logLevel)
        if err == nil {
            Logger.SetLevel(level)
        } else {
            Logger.Warn("Invalid LOG_LEVEL, defaulting to 'info'")
            Logger.SetLevel(logrus.InfoLevel)
        }
    } else {
        Logger.SetLevel(logrus.InfoLevel)
    }
}
