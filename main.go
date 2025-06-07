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
	"github.com/google/uuid"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	database := &db.SchipperkesDB{}
	database.Setup("data.db")
	admin := database.GetAdminUser()
	announcementList := database.GetAnnouncements()
	if slices.Contains(os.Args, "debug") {
		gin.SetMode(gin.DebugMode)
	}
	authManager := auth.AuthManager{Admin: &admin, DB: database}
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))

	store.Options(sessions.Options{MaxAge: int(30 * time.Minute), Path: "/", HttpOnly: true, Secure: true})
	router.Use(sessions.Sessions("admin-session", store))
	router.LoadHTMLGlob("frontend/**/*.html")
	router.Static("/styles", "./frontend/styles/")
	router.Static("/images", "./frontend/images/")
	router.Static("/scripts", "./frontend/scripts/")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.GET("/geschiedenis", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "geschiedenis.html", gin.H{})
	})

	router.GET("/admin/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "adminlogin.html", gin.H{})
	})

	adminGroup := router.Group("/admin")
	adminGroup.Use(authManager.AuthMiddleware())

	adminGroup.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "adminpanel.html", gin.H{})
	})

	adminGroup.GET("/announcements", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "adminannouncements.html", gin.H{})
	})

	adminGroup.POST("/login", authManager.LoginHandler)

	announcementApi := router.Group("/admin/announcements/api")
	announcementApi.GET("/get-all", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "announcementListTemplate.html", gin.H{
			"announcementList": announcementList,
		})
	})

	announcementApi.POST("/submit", func(ctx *gin.Context) {
		message := ctx.PostForm("message")
		newAnnouncement := db.Announcement{
			UUID:    uuid.New().String(),
			Message: message,
		}
		database.AddAnnouncement(&newAnnouncement)
		announcementList = database.GetAnnouncements()
		ctx.HTML(http.StatusOK, "announcementListTemplate.html", gin.H{
			"announcementList": announcementList,
		})
	})

	announcementApi.POST("/remove", func(ctx *gin.Context) {
		uuid := ctx.PostForm("UUID")
		database.RemoveAnnouncementByUUID(uuid)
		announcementList = database.GetAnnouncements()
	})
	router.Run(":8080")
}
