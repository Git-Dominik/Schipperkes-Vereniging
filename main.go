package main

import (
	"Git-Dominik/Schipperkes-Vereniging/auth"
	"Git-Dominik/Schipperkes-Vereniging/db"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func loadTemplates(pattern string) (*template.Template, error) {
	tmpl := template.New("")
	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			_, err := tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
		return
	}
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
	// router.LoadHTMLGlob("./**/*.html")
	tmpl, err := loadTemplates(".html")
	if err != nil {
		panic("Could not load templates: " + err.Error())
	}

	router.SetHTMLTemplate(tmpl)
	router.Static("/styles", "./frontend/styles/")
	router.Static("/images", "./frontend/images/")
	router.Static("/scripts", "./frontend/scripts/")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.GET("/geschiedenis", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "geschiedenis.html", gin.H{})
	})

	router.GET("/bestuur", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "bestuur.html", gin.H{})
	})

	router.GET("/admin/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "adminlogin.html", gin.H{})
	})

	router.POST("/admin/login", authManager.LoginHandler)

	adminGroup := router.Group("/admin")
	adminGroup.Use(authManager.AuthMiddleware())

	adminGroup.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "adminpanel.html", gin.H{})
	})

	adminGroup.GET("/announcements", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "adminannouncements.html", gin.H{})
	})

	announcementApi := router.Group("/admin/announcements/api")

	announcementApi.GET("/get/all", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "announcementListTemplate.html", gin.H{
			"announcementList": database.GetAnnouncements(),
		})
	})

	announcementApi.GET("/get/:uuid", func(ctx *gin.Context) {
		uuid := ctx.Param("uuid")
		announcement, err := database.GetAnnouncementByUUID(uuid)
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "", gin.H{})
		}
		// Uses list here but works in this case
		announcementList := []db.Announcement{*announcement}
		ctx.HTML(http.StatusOK, "announcementListTemplate.html", gin.H{"announcementList": announcementList})
	})

	announcementApi.GET("/get/:uuid/edit", func(ctx *gin.Context) {
		uuid := ctx.Param("uuid")
		announcement, err := database.GetAnnouncementByUUID(uuid)
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "", gin.H{})
		}
		ctx.HTML(http.StatusOK, "announcementEditTemplate.html", gin.H{"announcement": announcement})
	})

	announcementApi.PUT("/update/:uuid", func(ctx *gin.Context) {
		uuid := ctx.Param("uuid")
		title := ctx.PostForm("titel")
		message := ctx.PostForm("bericht")
		location := ctx.PostForm("locatie")

		announcement, err := database.GetAnnouncementByUUID(uuid)
		if err != nil {
			fmt.Println("not found")
			ctx.HTML(http.StatusBadRequest, "", gin.H{})
		}
		announcement.Message = message
		announcement.Title = title
		announcement.Location = location
		database.GormDB.Save(&announcement)
		// Uses list here but works in this case
		announcementList := []db.Announcement{*announcement}
		ctx.HTML(http.StatusOK, "announcementListTemplate.html", gin.H{"announcementList": announcementList})
	})

	announcementApi.POST("/submit", func(ctx *gin.Context) {
		title := ctx.PostForm("titel")
		message := ctx.PostForm("bericht")
		location := ctx.PostForm("locatie")
		newAnnouncement := db.Announcement{
			UUID:     uuid.New().String(),
			Message:  message,
			Title:    title,
			Location: location,
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
