package cmd

import (
	"fmt"
	"log"
	"sort"

	"github.com/adirkuhn/mig/migrations"
	"github.com/spf13/cobra"
)

type MigrationModel struct {
	ID string `gorm:"primaryKey"`
}

func (MigrationModel) TableName() string {
	return "migrations"
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply all pending migrations",
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

		for _, m := range migrations.Migrations {
			if !appliedMap[m.ID] {
				fmt.Println("Applying migration:", m.ID)
				if err := m.Up(db); err != nil {
					log.Fatalf("failed to apply migration %s: %v", m.ID, err)
				}
				db.Create(&MigrationModel{ID: m.ID})
			}
		}

		fmt.Println("Migrations applied successfully")
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}
