package cmd

import (
	"fmt"
	"log"
	"sort"

	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback the last applied migration",
	Run: func(cmd *cobra.Command, args []string) {

		var lastMigration MigrationModel
		if err := db.Order("id desc").First(&lastMigration).Error; err != nil {
			fmt.Println("No migrations to rollback")
			return
		}

		allMigrations := GetMigrations()

		sort.Slice(allMigrations, func(i, j int) bool {
			return allMigrations[i].ID > allMigrations[j].ID
		})

		for _, m := range allMigrations {
			if m.ID != lastMigration.ID {
				continue
			}

			fmt.Println("Rolling back migration:", m.ID)

			execDB := DB()

			if err := m.Down(execDB); err != nil {
				log.Fatalf("failed to rollback migration %s: %v", m.ID, err)
			}

			if dryRun {
				fmt.Println("-- dry-run: migration state not updated")
				return
			}

			db.Delete(&MigrationModel{ID: m.ID})
			fmt.Println("Migration rolled back successfully")
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(rollbackCmd)
}
