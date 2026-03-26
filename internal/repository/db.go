package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenPostgres(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&userModel{},
		&noteModel{},
		&markModel{},
		&sessionModel{},
	)
}
