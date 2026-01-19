package cmd

import (
	"fmt"
	"log"
	"sort"

	"github.com/adirkuhn/mig/migrations"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available migrations and their status",
	Run: func(cmd *cobra.Command, args []string) {
		err := db.AutoMigrate(&MigrationModel{})
		if err != nil {
			log.Fatal("failed to migrate migrations table")
		}

		var appliedMigrations []MigrationModel
		db.Find(&appliedMigrations)

		appliedMap := make(map[string]bool)
		for _, m := range appliedMigrations {
			appliedMap[m.ID] = true
		}

		sort.Slice(migrations.Migrations, func(i, j int) bool {
			return migrations.Migrations[i].ID < migrations.Migrations[j].ID
		})

		fmt.Println("Available migrations:")
		for _, m := range migrations.Migrations {
			status := "pending"
			if appliedMap[m.ID] {
				status = "applied"
			}
			fmt.Printf("  %s: %s\n", m.ID, status)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
