package db

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	GormDB gorm.DB
}

type Admin struct {
	gorm.Model
	HashedPassword []byte
	Email          string
}

type Announcement struct {
	UUID    uuid.UUID
	Message string
}

func (db *DB) Setup(databaseName string) {
	gormDB, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{})
	if err != nil {
		panic("Could not connect to database")
	}
	fmt.Println("Succesfully started database: " + databaseName)
	gormDB.AutoMigrate(&Admin{})
	gormDB.AutoMigrate(&Announcement{})
}
