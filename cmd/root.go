package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var db *gorm.DB

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

func NewMigratorCmd(gormDB *gorm.DB) *cobra.Command {
	db = gormDB
	return RootCmd
}
