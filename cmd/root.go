package cmd

import (
    "github.com/spf13/cobra"
    "github.com/sirupsen/logrus"
)

var rootCmd = &cobra.Command{
    Use:   "trail-cli",
    Short: "Trail Finder CLI",
    Long:  "CLI application for managing trail data",
}

// Execute runs the root command
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        logrus.Fatalf("Error executing command: %v", err)
    }
}
