package cmd

import (
	"fmt"
	"log"
	"sort"

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

		allMigrations := GetMigrations() // Updated to GetMigrations()

		sort.Slice(allMigrations, func(i, j int) bool {
			return allMigrations[i].ID < allMigrations[j].ID
		})

		fmt.Println("Available migrations:")
		for _, m := range allMigrations {
			status := "pending"
			if appliedMap[m.ID] {
				status = "applied"
			}
			fmt.Printf("  %s [%s]: %s\n", m.ID, m.Name, status) // Updated to include m.Name
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
