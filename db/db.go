package db

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type SchipperkesDB struct {
	GormDB *gorm.DB
}

type Admin struct {
	gorm.Model
	HashedPassword []byte
	Email          string
	AdminUUID      string
}

func (schipperkesDB *SchipperkesDB) Setup(databaseName string) {
	db, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("Could not connect to database")
	}
	fmt.Println("Succesfully started database: " + databaseName)
	db.AutoMigrate(&Admin{})
	db.AutoMigrate(&Activity{})
	db.AutoMigrate(&SignUp{})
	schipperkesDB.GormDB = db
}
