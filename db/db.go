package db

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
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

type Announcement struct {
	UUID    uuid.UUID
	Message string
}

func (schipperkesDB *SchipperkesDB) Setup(databaseName string) {
	db, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("Could not connect to database")
	}
	fmt.Println("Succesfully started database: " + databaseName)
	db.AutoMigrate(&Admin{})
	db.AutoMigrate(&Announcement{})
	schipperkesDB.GormDB = db
}

func (schipperkesDB *SchipperkesDB) GetAdminUser() Admin {
	db := schipperkesDB.GormDB
	var admin Admin
	err := db.First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("No admin account found creating default account.")
		admin.Email = "ezraschutte227@gmail.com"
		defaultPassword := "admin"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
			panic("Error creating admin account.")
		}
		admin.HashedPassword = hashedPassword
		db.Create(&admin)
	}

	return admin
}
