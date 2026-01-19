package cmd

import (
	"fmt"
	"log"
	"sort"

	"github.com/adirkuhn/mig/migrations"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var dryrunCmd = &cobra.Command{
	Use:   "dry-run",
	Short: "Show the SQL statements for pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		appliedDb := db.Session(&gorm.Session{})

		err := appliedDb.AutoMigrate(&MigrationModel{})
		if err != nil {
			log.Fatal("failed to migrate migrations table")
		}

		var appliedMigrations []MigrationModel
		appliedDb.Find(&appliedMigrations)

		appliedMap := make(map[string]bool)
		for _, m := range appliedMigrations {
			appliedMap[m.ID] = true
		}

		sort.Slice(migrations.Migrations, func(i, j int) bool {
			return migrations.Migrations[i].ID < migrations.Migrations[j].ID
		})

		fmt.Println("Pending migrations SQL:")
		for _, m := range migrations.Migrations {
			if !appliedMap[m.ID] {
				fmt.Println("-- Migration:", m.ID)
				dryRunDb := db.Session(&gorm.Session{DryRun: true})
				if err := m.Up(dryRunDb); err != nil {
					log.Fatalf("failed to generate SQL for migration %s: %v", m.ID, err)
				}
				if dryRunDb.Statement.SQL.String() != "" {
					fmt.Println(dryRunDb.Statement.SQL.String())
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(dryrunCmd)
}
