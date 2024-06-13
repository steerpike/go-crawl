package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var createDbCmd = &cobra.Command{
	Use:   "createdb",
	Short: "Create a SQLite database from a schema.sql file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbName := args[0]
		cmdStr := fmt.Sprintf("sqlite3 %s < db/schema.sql", dbName)
		out, err := exec.Command("sh", "-c", cmdStr).Output()
		if err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		fmt.Printf("Output: %s\n", out)
	},
}

var dropDbCmd = &cobra.Command{
	Use:   "dropdb",
	Short: "Drop a SQLite database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbName := args[0]
		err := os.Remove(dbName)
		if err != nil {
			log.Fatalf("Failed to drop database: %v", err)
		} else {
			fmt.Printf("Database %s dropped successfully\n", dbName)
		}
	},
}

func init() {
	rootCmd.AddCommand(createDbCmd)
	rootCmd.AddCommand(dropDbCmd)
}
