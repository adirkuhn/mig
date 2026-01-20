package cmd

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
SQLCaptureLogger captures ALL SQL statements executed by GORM
during DryRun mode. This is the ONLY reliable way to dump SQL.
*/
type SQLCaptureLogger struct{}

func (l *SQLCaptureLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *SQLCaptureLogger) Info(ctx context.Context, msg string, data ...interface{})  {}
func (l *SQLCaptureLogger) Warn(ctx context.Context, msg string, data ...interface{})  {}
func (l *SQLCaptureLogger) Error(ctx context.Context, msg string, data ...interface{}) {}

func (l *SQLCaptureLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (string, int64),
	err error,
) {
	sql, _ := fc()
	if sql != "" {
		fmt.Println(sql + ";")
	}
}

var dryrunCmd = &cobra.Command{
	Use:   "dry-run",
	Short: "Show the SQL statements for pending migrations",
	Run: func(cmd *cobra.Command, args []string) {

		// Ensure migrations table exists (real DB, not dry-run)
		if err := db.AutoMigrate(&MigrationModel{}); err != nil {
			log.Fatal("failed to migrate migrations table")
		}

		// Load applied migrations
		var appliedMigrations []MigrationModel
		if err := db.Find(&appliedMigrations).Error; err != nil {
			log.Fatal("failed to load applied migrations")
		}

		appliedMap := make(map[string]bool)
		for _, m := range appliedMigrations {
			appliedMap[m.ID] = true
		}

		// Load all registered migrations
		allMigrations := GetMigrations()

		sort.Slice(allMigrations, func(i, j int) bool {
			return allMigrations[i].ID < allMigrations[j].ID
		})

		fmt.Println("Pending migrations SQL:")

		for _, m := range allMigrations {
			if appliedMap[m.ID] {
				continue
			}

			fmt.Printf("\n-- Migration: %s (%s)\n", m.ID, m.Name)

			// DryRun DB with SQL capture logger
			dryRunDb := db.Session(&gorm.Session{
				DryRun: true,
				Logger: &SQLCaptureLogger{},
			})

			// Execute migration logic (NO DB writes)
			if err := m.Up(dryRunDb); err != nil {
				log.Fatalf(
					"failed to generate SQL for migration %s: %v",
					m.ID,
					err,
				)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(dryrunCmd)
}
