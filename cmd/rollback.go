package cmd

import (
	"fmt"
	"log"
	"sort"

	"github.com/adirkuhn/mig/migrations"
	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback the last applied migration",
	Run: func(cmd *cobra.Command, args []string) {
		err := db.AutoMigrate(&MigrationModel{})
		if err != nil {
			log.Fatal("failed to migrate migrations table")
		}

		var lastMigration MigrationModel
		db.Order("id desc").First(&lastMigration)

		if lastMigration.ID == "" {
			fmt.Println("No migrations to rollback")
			return
		}

		sort.Slice(migrations.Migrations, func(i, j int) bool {
			return migrations.Migrations[i].ID > migrations.Migrations[j].ID
		})

		for _, m := range migrations.Migrations {
			if m.ID == lastMigration.ID {
				fmt.Println("Rolling back migration:", m.ID)
				if err := m.Down(db); err != nil {
					log.Fatalf("failed to rollback migration %s: %v", m.ID, err)
				}
				db.Delete(&MigrationModel{ID: m.ID})
				fmt.Println("Migration rolled back successfully")
				return
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(rollbackCmd)
}
