package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setup(t *testing.T) (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	return db, func() {
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		sqlDB.Close()
	}
}

func execute(t *testing.T, cmd *cobra.Command, args ...string) (string, error) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	err := cmd.Execute()

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	return string(out) + buf.String(), err
}

// filterGormLog removes common GORM log prefixes from the output.
func filterGormLog(output string) string {
	var filteredLines []string
	lines := bytes.Split([]byte(output), []byte("\n"))
	for _, line := range lines {
		if !bytes.Contains(line, []byte("gorm.io/gorm")) && !bytes.Contains(line, []byte("[0.000ms]")) {
			filteredLines = append(filteredLines, string(line))
		}
	}
	return strings.Join(filteredLines, "\n")
}

