package main

import (
	"Git-Dominik/Schipperkes-Vereniging/auth"
	"Git-Dominik/Schipperkes-Vereniging/db"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	db := &db.SchipperkesDB{}
	db.Setup("data.db")
	admin := db.GetAdminUser()
	if slices.Contains(os.Args, "debug") {
		gin.SetMode(gin.DebugMode)
	}
	authManager := auth.AuthManager{Admin: &admin, DB: db}
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))

	store.Options(sessions.Options{MaxAge: int(30 * time.Minute), Path: "/", HttpOnly: true, Secure: true}) // Set session options
	router.Use(sessions.Sessions("admin-session", store))
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

	adminGroup := router.Group("/admin")
	adminGroup.Use(authManager.AuthMiddleware())

	adminGroup.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "adminpanel.html", gin.H{})
	})

	router.GET("/admin/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "adminlogin.html", gin.H{})
	})

	adminGroup.POST("/login", authManager.LoginHandler)

	router.Run(":8080")
}
