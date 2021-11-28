package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectSQLite(dbFile string) (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		log.Fatalf("Connect SQLite: %s\n", err)
	}
	return
}
