package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	if slices.Contains(os.Args, "debug") {
		gin.SetMode(gin.DebugMode)
	}
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
		ctx.HTML(http.StatusOK, "adminlogin.html", gin.H{})
	})

	router.POST("/admin/login", func(ctx *gin.Context) {
		password := ctx.PostForm("adminPassword")
		// Type = []uint8
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
			ctx.HTML(http.StatusOK, "loginfailed.html", gin.H{})
		}
		fmt.Println(hashedPassword)

	})

	router.Run(":8080")
}
