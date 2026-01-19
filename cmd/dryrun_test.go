package cmd

import (
	"testing"

	"github.com/adirkuhn/mig/migrations"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type Product struct {
	ID    uint   `gorm:"primaryKey"`
	Code  string `gorm:"size:255"`
	Price uint
}

func TestDryRunCmd(t *testing.T) {
	db, teardown := setup(t)
	defer teardown()

	cmd := NewMigratorCmd(db)

	// Create dummy migrations
	originalMigrations := migrations.Migrations
	migrations.Migrations = []*migrations.Migration{
		{
			ID:   "20240101120000",
			Up:   func(db *gorm.DB) error { return db.AutoMigrate(&Product{}) },
			Down: func(db *gorm.DB) error { return db.Migrator().DropTable(&Product{}) },
		},
	}
	defer func() { migrations.Migrations = originalMigrations }()

	// Run the dry-run command
	output, err := execute(t, cmd, "dry-run")
	assert.NoError(t, err)

	// Check if the output contains the SQL statements
	assert.Contains(t, output, "CREATE TABLE `products` (`id` integer PRIMARY KEY AUTOINCREMENT,`code` text,`price` integer)")
	assert.NotContains(t, output, "INSERT INTO `migrations`") // Should not apply migration

	// Verify that the table was not created
	assert.False(t, db.Migrator().HasTable("products"))
}
