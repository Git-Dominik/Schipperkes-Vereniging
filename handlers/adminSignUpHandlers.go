package handlers

import (
	"Git-Dominik/Schipperkes-Vereniging/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllSignUps(database *db.SchipperkesDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		signUpList := database.GetAllSignUps()
		ctx.HTML(http.StatusOK, "adminSignUpList.html", gin.H{
			"signUpList": signUpList,
		})
	}
}

func RemoveSignUp(database *db.SchipperkesDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuid := ctx.PostForm("UUID")
		database.RemoveSignUpByUUID(uuid)
	}
}
