package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var db *gorm.DB
var dryRun bool

var RootCmd = &cobra.Command{
	Use:   "migrator",
	Short: "A simple migration tool for GORM",
	Long:  `A simple migration tool for GORM to manage database migrations.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func SetDB(gormDB *gorm.DB) {
	db = gormDB
}

func NewMigratorCmd(gormDB *gorm.DB) *cobra.Command {
	SetDB(gormDB)
	// Migration schema
	if err := db.AutoMigrate(&MigrationModel{}); err != nil {
		log.Fatal("failed to migrate migrations table")
	}
	return RootCmd
}

func init() {
	RootCmd.PersistentFlags().BoolVar(
		&dryRun,
		"dry-run",
		false,
		"Print SQL without executing migrations",
	)
}
