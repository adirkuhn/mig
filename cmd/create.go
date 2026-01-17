package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var migrationsDir string

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		timestamp := time.Now().Format("20060102150405")
		filename := fmt.Sprintf("%s_%s.go", timestamp, name)
		
		if migrationsDir == "" {
			migrationsDir = "migrations" // Default to "migrations" if not set
		}

		// Ensure the migrations directory exists
		if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
			if err = os.MkdirAll(migrationsDir, os.ModePerm); err != nil {
				fmt.Printf("Error creating migrations directory %s: %v\n", migrationsDir, err)
				return
			}
		}
		
		filepath := filepath.Join(migrationsDir, filename)

		file, err := os.Create(filepath)
		if err != nil {
			fmt.Println("Error creating migration file:", err)
			return
		}
		defer file.Close()

		template := fmt.Sprintf(`package migrations

import (
	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "%s",
		Up: func(db *gorm.DB) error {
			// your migration logic here
			return nil
		},
		Down: func(db *gorm.DB) error {
			// your rollback logic here
			return nil
		},
	})
}
`, timestamp)

		_, err = file.WriteString(template)
		if err != nil {
			fmt.Println("Error writing to migration file:", err)
			return
		}

		fmt.Println("Created migration:", filepath)
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&migrationsDir, "dir", "d", "", "Directory to store migration files")
}
