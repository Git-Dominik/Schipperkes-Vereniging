package db

import (
	"errors"
	"fmt"
	"os"

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
	gorm.Model
	Title    string
	UUID     string
	Message  string
	Location string
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

func (schipperkesDB *SchipperkesDB) GetAnnouncementByUUID(uuid string) (*Announcement, error) {
	db := schipperkesDB.GormDB
	var announcement Announcement
	if err := db.Where("UUID = ?", uuid).First(&announcement).Error; err != nil {
		return nil, err
	}
	return &announcement, nil
}

func (schipperkesDB *SchipperkesDB) GetAdminUser() Admin {
	db := schipperkesDB.GormDB
	var admin Admin
	err := db.First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("No admin account found creating default account.")
		admin.Email = os.Getenv("DEFAULT_EMAIL")
		defaultPassword := os.Getenv("DEFAULT_PASSWORD")
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

func (schipperkesDB *SchipperkesDB) GetAnnouncements() []Announcement {
	db := schipperkesDB.GormDB
	var announcements []Announcement
	err := db.Find(&announcements).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		announcements = []Announcement{}
	}
	return announcements
}

func (schipperkesDB *SchipperkesDB) AddAnnouncement(announcement *Announcement) {
	db := schipperkesDB.GormDB
	db.Create(announcement)
}

func (schipperkesDB *SchipperkesDB) RemoveAnnouncementByUUID(uuid string) {
	db := schipperkesDB.GormDB
	db.Where("UUID = ?", uuid).Delete(&Announcement{})
}
