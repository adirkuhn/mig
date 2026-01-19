# Migrator

A simple database migration tool for GORM.

## Usage

The migrator tool provides several commands to manage your database migrations.

### Create a Migration

To create a new migration file, use the `create` command:

```bash
migrator create <migration_name> --dir <migrations_directory>
```

This will create a new migration file with the current timestamp and the provided name in the specified directory.

### Apply Migrations

To apply all pending migrations, use the `migrate` command:

```bash
migrator migrate
```

### Rollback Migrations

To roll back the last applied migration, use the `rollback` command:

```bash
migrator rollback
```

### Dry Run

To see which migrations are pending without applying them, use the `dryrun` command:

```bash
migrator dryrun
```

### List Migrations

To see all applied migrations, use the `list` command:

```bash
migrator list
```

### Set Database Connection

To set the database connection string, use the `set` command:

```bash
migrator set <database_url>
```

## How to use it in your project

You can integrate the migrator with your Go application by using the `NewMigratorCmd` function from the `cmd` package.

Here is an example of how to use it:

```go
package main

import (
	"log"
	"migrator/cmd"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Initialize your GORM database connection
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	// Create a new migrator command
	migratorCmd := cmd.NewMigratorCmd(db)

	// Execute the command
	if err := migratorCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
```

Then you can run the migrator from your application:

```bash
go run main.go migrate
```

You can also set the `MIGRATIONS_DIR` environment variable to specify the directory where your migration files are stored.
