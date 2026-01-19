package cmd

import (
	"fmt"
	"log"
	"sort"

	"github.com/adirkuhn/mig/migrations"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set [version]",
	Short: "Migrate or rollback to a specific version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]

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

		for _, m := range migrations.Migrations {
			if m.ID > version {
				if appliedMap[m.ID] {
					fmt.Println("Rolling back migration:", m.ID)
					if err := m.Down(db); err != nil {
						log.Fatalf("failed to rollback migration %s: %v", m.ID, err)
					}
					db.Delete(&MigrationModel{ID: m.ID})
				}
			}
		}

		sort.Slice(migrations.Migrations, func(i, j int) bool {
			return migrations.Migrations[i].ID < migrations.Migrations[j].ID
		})

		for _, m := range migrations.Migrations {
			if m.ID <= version {
				if !appliedMap[m.ID] {
					fmt.Println("Applying migration:", m.ID)
					if err := m.Up(db); err != nil {
						log.Fatalf("failed to apply migration %s: %v", m.ID, err)
					}
					db.Create(&MigrationModel{ID: m.ID})
				}
			}
		}

		fmt.Println("Migrations set to version", version)
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
