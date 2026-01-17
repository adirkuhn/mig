package migrations

import "gorm.io/gorm"

type Migration struct {
	ID   string
	Up   func(db *gorm.DB) error
	Down func(db *gorm.DB) error
}

var Migrations []*Migration
