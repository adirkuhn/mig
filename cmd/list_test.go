package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestListCmd(t *testing.T) {
	db, teardown := setup(t)
	defer teardown()

	cmd := NewMigratorCmd(db)

	// Create dummy migrations
	Register(&Migration{
		ID:   "20240101120000",
		Name: "CreateDummyTable",
		Up:   func(db *gorm.DB) error { return nil },
		Down: func(db *gorm.DB) error { return nil },
	})
	Register(&Migration{
		ID:   "20240101120001",
		Name: "AddAnotherDummyTable",
		Up:   func(db *gorm.DB) error { db.Exec("CREATE TABLE dummy (id INT)"); return nil },
		Down: func(db *gorm.DB) error { db.Exec("DROP TABLE dummy"); return nil },
	})
	defer ClearRegistry()

	// List without applying any migrations
	output, err := execute(t, cmd, "list")
	assert.NoError(t, err)
	filteredOutput := filterGormLog(output)
	assert.Contains(t, filteredOutput, "20240101120000 [CreateDummyTable]: pending")
	assert.Contains(t, filteredOutput, "20240101120001 [AddAnotherDummyTable]: pending")

	// Apply one migration
	_, err = execute(t, cmd, "migrate")
	assert.NoError(t, err)

	// List after applying one migration
	output, err = execute(t, cmd, "list")
	assert.NoError(t, err)
	filteredOutput = filterGormLog(output)
	assert.Contains(t, filteredOutput, "20240101120000 [CreateDummyTable]: applied")
	assert.Contains(t, filteredOutput, "20240101120001 [AddAnotherDummyTable]: applied")
}
