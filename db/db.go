package db

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
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

type Activity struct {
	gorm.Model
	Title         string
	UUID          string
	Message       string
	Location      string
	DateTimeStart time.Time
	DateTimeEnd   time.Time
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
	schipperkesDB.GormDB = db
}

func (schipperkesDB *SchipperkesDB) GetActivityByUUID(uuid string) (*Activity, error) {
	db := schipperkesDB.GormDB
	var activity Activity
	if err := db.Where("UUID = ?", uuid).First(&activity).Error; err != nil {
		return nil, err
	}
	return &activity, nil
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

func (schipperkesDB *SchipperkesDB) GetActivities() []Activity {
	db := schipperkesDB.GormDB
	var activityList []Activity
	err := db.Find(&activityList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		activityList = []Activity{}
	}
	return activityList
}

func (schipperkesDB *SchipperkesDB) AddActivity(activity *Activity) {
	db := schipperkesDB.GormDB
	db.Create(activity)
}

func (schipperkesDB *SchipperkesDB) RemoveActivityByUUID(uuid string) {
	db := schipperkesDB.GormDB
	db.Where("UUID = ?", uuid).Delete(&Activity{})
}
