package cmd

import (
	"testing"

	"github.com/adirkuhn/mig/migrations"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRollbackCmd(t *testing.T) {
	t.Run("rolls back the last applied migration", func(t *testing.T) {
		db, teardown := setup(t)
		defer teardown()

		cmd := NewMigratorCmd(db)

		// Create a dummy migration
		originalMigrations := migrations.Migrations
		migrations.Migrations = []*migrations.Migration{
			{
				ID: "20240101120000",
				Up: func(db *gorm.DB) error {
					return db.Exec("CREATE TABLE test_table (id INT)").Error
				},
				Down: func(db *gorm.DB) error {
					return db.Exec("DROP TABLE test_table").Error
				},
			},
		}
		defer func() { migrations.Migrations = originalMigrations }()

		// Apply the migration
		_, err := execute(t, cmd, "migrate")
		assert.NoError(t, err)

		// Run the rollback command
		output, err := execute(t, cmd, "rollback")
		assert.NoError(t, err)
		filteredOutput := filterGormLog(output)
		assert.Contains(t, filteredOutput, "Rolling back migration: 20240101120000")
		assert.Contains(t, filteredOutput, "Migration rolled back successfully")

		// Check if the migration was rolled back
		var count int64
		db.Raw("SELECT count(*) FROM migrations WHERE id = ?", "20240101120000").Scan(&count)
		assert.Equal(t, int64(0), count)

		// Check if the table was dropped
		assert.False(t, db.Migrator().HasTable("test_table"))
	})
}
