package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateCmd(t *testing.T) {
	db, teardown := setup(t)
	defer teardown()

	cmd := NewMigratorCmd(db)

	// Create a temporary directory for migrations
	tempDir, err := os.MkdirTemp("", "migrations_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir) // Clean up the temporary directory

	// Test creating a new migration
	name := "test_migration"
	output, err := execute(t, cmd, "create", name, "--dir", tempDir)
	assert.NoError(t, err)

	timestamp := time.Now().Format("20060102150405")
	generatedFilename := fmt.Sprintf("%s_%s.go", timestamp, name)
	expectedOutputSubstring := fmt.Sprintf("Created migration: %s", filepath.Join(tempDir, generatedFilename))
	assert.Contains(t, output, expectedOutputSubstring)

	// Check if the file was created
	_, err = os.Stat(filepath.Join(tempDir, generatedFilename))
	assert.NoError(t, err)
}
