package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/adirkuhn/mig/cmd"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	if os.Getenv("MIGRATIONS_DIR") == "" {
		log.Fatal("MIGRATIONS_DIR environment variable not set")
	}

	parts := strings.SplitN(databaseURL, "://", 2)
	if len(parts) != 2 {
		log.Fatalf("invalid DATABASE_URL format: %s", databaseURL)
	}

	driver := parts[0]
	dsn := parts[1]
	var db *gorm.DB

	switch driver {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	case "postgres", "postgresql":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	default:
		log.Fatalf("unsupported db driver: %s", driver)
	}

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	var rootCmd = &cobra.Command{
		Use:   "my-app",
		Short: "My application with migrations",
	}

	migratorCmd := cmd.NewMigratorCmd(db)
	rootCmd.AddCommand(migratorCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
