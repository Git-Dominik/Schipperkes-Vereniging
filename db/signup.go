package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type SignUp struct {
	gorm.Model
	UUID            string
	FirstName       string
	LastName        string
	Email           string
	Phone           string
	Birthdate       time.Time
	StreetAndNumber string
	ZipCode         string
	City            string
	ExtraInfo       string
}

func (schipperkesDB *SchipperkesDB) AddSignUp(signUp *SignUp) {
	gormDB := schipperkesDB.GormDB
	gormDB.Create(signUp)
}

func (schipperkesDB *SchipperkesDB) RemoveSignUpByUUID(uuid string) {
	db := schipperkesDB.GormDB
	db.Where("UUID = ?", uuid).Delete(&SignUp{})
}
func (schipperkesDB *SchipperkesDB) GetAllSignUps() []SignUp {
	db := schipperkesDB.GormDB
	var signUpList []SignUp
	err := db.Find(&signUpList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		signUpList = []SignUp{}
	}
	return signUpList
}
