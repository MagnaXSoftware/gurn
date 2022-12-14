package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabases(dsn string) error {
	// Just init default for now

	newDb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return err
	}

	err = newDb.AutoMigrate(&Urn{}, &Access{})
	if err != nil {
		return err
	}

	db = newDb

	return nil
}

func GetDatabase() *gorm.DB {
	return db
}
