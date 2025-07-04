package handlers

import (
	"Git-Dominik/Schipperkes-Vereniging/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ContactForm(database *db.SchipperkesDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuid := uuid.New().String()
		firstName := ctx.PostForm("firstName")
		lastName := ctx.PostForm("lastName")
		email := ctx.PostForm("email")
		phone := ctx.PostForm("phone")
		streetAndNumber := ctx.PostForm("street")
		zipCode := ctx.PostForm("zipCode")
		city := ctx.PostForm("city")
		extraInfo := ctx.PostForm("extraInfo")

		birthdate := ctx.PostForm("birthDate")
		birthdataParsed, err := time.Parse("2006-01-02", birthdate)
		if err != nil {
			ctx.HTML(http.StatusOK, "failSignUp.html", gin.H{})
			return
		}

		signUp := db.SignUp{
			UUID:            uuid,
			FirstName:       firstName,
			LastName:        lastName,
			Email:           email,
			Phone:           phone,
			Birthdate:       birthdataParsed,
			StreetAndNumber: streetAndNumber,
			ZipCode:         zipCode,
			City:            city,
			ExtraInfo:       extraInfo,
		}
		database.AddSignUp(&signUp)
		ctx.HTML(http.StatusOK, "succesfullSignUp.html", gin.H{})
	}
}
