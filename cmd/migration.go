package cmd

import (
	"log"

	"gorm.io/gorm"
)

type Migration struct {
	ID   string
	Name string
	Up   func(db *gorm.DB) error
	Down func(db *gorm.DB) error
}

var migrationRegistry = make([]*Migration, 0)

func Register(m *Migration) {
	migrationRegistry = append(migrationRegistry, m)
}

func GetMigrations() []*Migration {
	log.Printf("Getting migrations. Registry contains %d migrations.", len(migrationRegistry)) // Added for debugging
	for _, m := range migrationRegistry {
		log.Printf("  Migration ID: %s [%s]", m.ID, m.Name) // Added for debugging
	}
	return migrationRegistry
}

func ClearRegistry() {
	migrationRegistry = make([]*Migration, 0)
}
