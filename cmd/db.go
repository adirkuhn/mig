package cmd

import "gorm.io/gorm"

// DB returns the correct DB session depending on --dry-run
func DB() *gorm.DB {
	if dryRun {
		return db.Session(&gorm.Session{
			DryRun: true,
			Logger: &SQLCaptureLogger{},
		})
	}

	return db
}
