package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSetCmd(t *testing.T) {
	db, teardown := setup(t)
	defer teardown()

	cmd := NewMigratorCmd(db)

	// Create dummy migrations
	Register(&Migration{
		ID:   "20240101120000",
		Up:   func(db *gorm.DB) error { return db.Exec("CREATE TABLE table1 (id INT)").Error },
		Down: func(db *gorm.DB) error { return db.Exec("DROP TABLE table1").Error },
	})
	Register(&Migration{
		ID:   "20240101120001",
		Up:   func(db *gorm.DB) error { return db.Exec("CREATE TABLE table2 (id INT)").Error },
		Down: func(db *gorm.DB) error { return db.Exec("DROP TABLE table2").Error },
	})
	Register(&Migration{
		ID:   "20240101120002",
		Up:   func(db *gorm.DB) error { return db.Exec("CREATE TABLE table3 (id INT)").Error },
		Down: func(db *gorm.DB) error { return db.Exec("DROP TABLE table3").Error },
	})
	defer ClearRegistry()

	// Test migrating up to a specific version
	output, err := execute(t, cmd, "set", "20240101120001")
	assert.NoError(t, err)
	filteredOutput := filterGormLog(output)
	assert.Contains(t, filteredOutput, "Applying migration: 20240101120000")
	assert.Contains(t, filteredOutput, "Applying migration: 20240101120001")
	assert.Contains(t, filteredOutput, "Migrations set to version 20240101120001")

	assert.True(t, db.Migrator().HasTable("table1"))
	assert.True(t, db.Migrator().HasTable("table2"))
	assert.False(t, db.Migrator().HasTable("table3"))

	// Test migrating down to a specific version
	output, err = execute(t, cmd, "set", "20240101120000")
	assert.NoError(t, err)
	filteredOutput = filterGormLog(output)
	assert.Contains(t, filteredOutput, "Rolling back migration: 20240101120001")
	assert.Contains(t, filteredOutput, "Migrations set to version 20240101120000")
}
