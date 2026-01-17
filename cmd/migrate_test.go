package cmd

import (
	"migrator/migrations"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMigrateCmd(t *testing.T) {
	t.Run("applies pending migrations", func(t *testing.T) {
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

		// Run the migrate command
		output, err := execute(t, cmd, "migrate")
		assert.NoError(t, err)
		filteredOutput := filterGormLog(output)
		assert.Contains(t, filteredOutput, "Applying migration: 20240101120000")
		assert.Contains(t, filteredOutput, "Migrations applied successfully")

		// Check if the migration was applied
		var count int64
		db.Raw("SELECT count(*) FROM migrations WHERE id = ?", "20240101120000").Scan(&count)
		assert.Equal(t, int64(1), count)

		// Check if the table was created
		assert.True(t, db.Migrator().HasTable("test_table"))
	})
}
