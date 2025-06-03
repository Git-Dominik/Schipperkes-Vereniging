package main

import (
	"Git-Dominik/Schipperkes-Vereniging/db"
	"fmt"
	"net/http"
	"os"
	"slices"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	db := &db.SchipperkesDB{}
	db.Setup("data.db")
	admin := db.GetAdminUser()
	if slices.Contains(os.Args, "debug") {
		gin.SetMode(gin.DebugMode)
	}
	fmt.Println("Email:")
	fmt.Println(admin.Email)
	router := gin.Default()
	router.LoadHTMLGlob("frontend/*.html")
	router.Static("/styles", "./frontend/styles/")
	router.Static("/images", "./frontend/images/")
	router.Static("/scripts", "./frontend/scripts/")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.GET("/geschiedenis", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "geschiedenis.html", gin.H{})
	})

	router.GET("/admin", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "adminpanel.html", gin.H{})
	})

	router.GET("/admin-login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "adminlogin.html", gin.H{})
	})

	router.POST("/admin/login", func(ctx *gin.Context) {
		email := ctx.PostForm("adminEmail")
		password := ctx.PostForm("adminPassword")
		err := bcrypt.CompareHashAndPassword(admin.HashedPassword, []byte(password))
		if err != nil || email != admin.Email {
			ctx.HTML(http.StatusOK, "loginfailed.html", gin.H{})
			return
		}
		// Tell htmx to go to /admin
		ctx.Header("HX-Redirect", "/admin")
	})

	router.Run(":8080")
}
