package cmd

import (
	"migrator/migrations"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestListCmd(t *testing.T) {
	db, teardown := setup(t)
	defer teardown()

	cmd := NewMigratorCmd(db)

	// Create dummy migrations
	originalMigrations := migrations.Migrations
	migrations.Migrations = []*migrations.Migration{
		{
			ID: "20240101120000",
			Up: func(db *gorm.DB) error { return nil },
			Down: func(db *gorm.DB) error { return nil },
		},
		{
			ID: "20240101120001",
			Up: func(db *gorm.DB) error { return nil },
			Down: func(db *gorm.DB) error { return nil },
		},
	}
	defer func() { migrations.Migrations = originalMigrations }()

	// List without applying any migrations
	output, err := execute(t, cmd, "list")
	assert.NoError(t, err)
	filteredOutput := filterGormLog(output)
	assert.Contains(t, filteredOutput, "20240101120000: pending")
	assert.Contains(t, filteredOutput, "20240101120001: pending")

	// Apply one migration
	_, err = execute(t, cmd, "migrate")
	assert.NoError(t, err)

	// List after applying one migration
	output, err = execute(t, cmd, "list")
	assert.NoError(t, err)
	filteredOutput = filterGormLog(output)
	assert.Contains(t, filteredOutput, "20240101120000: applied")
	assert.Contains(t, filteredOutput, "20240101120001: applied")
}