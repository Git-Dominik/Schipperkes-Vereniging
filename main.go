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
	"strconv"
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
	database.Setup(filepath.Join("data" + "data.db"))
	admin := database.GetAdminUser()
	activityList := database.GetActivities()
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
	router.Static("/files", "./frontend/files/")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.GET("/geschiedenis", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "geschiedenis.html", gin.H{})
	})

	router.GET("/bestuur", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "bestuur.html", gin.H{})
	})

	router.GET("/activiteiten", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "activiteiten.html", gin.H{})
	})

	router.GET("/activiteiten/get/all", func(ctx *gin.Context) {
		// Make it easy to print in html without needing to process the data in javascript
		type easyFormat struct {
			DayEasyFormat       int
			MonthEasyFormat     string
			StartTimeEasyFormat string
			EndTimeEasyFormat   string
			Title               string
			Message             string
			Location            string
		}

		activitiesList := database.GetActivities()
		if len(activitiesList) <= 0 {
			ctx.HTML(http.StatusOK, "geenActiviteitenGeplandTemplate.html", gin.H{})
		}

		// Parse the optional "limit" query parameter
		limitParam := ctx.DefaultQuery("limit", "")
		var limit int
		var err error
		if limitParam != "" {
			limit, err = strconv.Atoi(limitParam)
			if err != nil || limit < 0 {
				ctx.String(http.StatusBadRequest, "Invalid limit parameter")
				return
			}
			if limit > len(activitiesList) {
				limit = len(activitiesList)
			}
			activitiesList = activitiesList[:limit]
		}

		listOfEasyFormat := []easyFormat{}
		for _, activity := range activitiesList {
			easyFormatItem := easyFormat{
				DayEasyFormat:       activity.DateTimeStart.Day(),
				MonthEasyFormat:     activity.DateTimeStart.Month().String(),
				StartTimeEasyFormat: activity.DateTimeStart.Format("15:04"),
				EndTimeEasyFormat:   activity.DateTimeEnd.Format("15:04"),
				Title:               activity.Title,
				Message:             activity.Message,
				Location:            activity.Location,
			}
			listOfEasyFormat = append(listOfEasyFormat, easyFormatItem)
		}

		ctx.HTML(http.StatusOK, "activityListTemplate.html", gin.H{
			"activityList": listOfEasyFormat,
		})
	})

	router.GET("/admin/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "adminLogin.html", gin.H{})
	})

	router.POST("/admin/login", authManager.LoginHandler)

	adminGroup := router.Group("/admin")
	adminGroup.Use(authManager.AuthMiddleware())

	adminGroup.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "adminPanel.html", gin.H{})
	})

	adminGroup.GET("/activities", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "adminActivities.html", gin.H{})
	})

	activityApi := router.Group("/admin/activities/api")

	activityApi.GET("/get/all", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "activityListTemplateAdmin.html", gin.H{
			"activityList": database.GetActivities(),
		})
	})

	activityApi.GET("/get/:uuid", func(ctx *gin.Context) {
		uuid := ctx.Param("uuid")
		activity, err := database.GetActivityByUUID(uuid)
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "", gin.H{})
		}
		// Uses list here but works in this case
		activityList := []db.Activity{*activity}
		ctx.HTML(http.StatusOK, "activityListTemplateAdmin.html", gin.H{"activityList": activityList})
	})

	activityApi.GET("/get/:uuid/edit", func(ctx *gin.Context) {
		uuid := ctx.Param("uuid")
		activity, err := database.GetActivityByUUID(uuid)
		if err != nil {
			ctx.String(http.StatusBadRequest, "Activity not found", gin.H{})
		}
		ctx.HTML(http.StatusOK, "activityEditTemplate.html", gin.H{"activity": activity})
	})

	activityApi.PUT("/update/:uuid", func(ctx *gin.Context) {
		uuid := ctx.Param("uuid")
		title := ctx.PostForm("titel")
		message := ctx.PostForm("bericht")
		location := ctx.PostForm("locatie")
		dateTimeStartPost := ctx.PostForm("datumTijdStart")
		dateTimeEindPost := ctx.PostForm("datumTijdEnd")
		dateTimeStart, err := time.Parse("2006-01-02T15:04", dateTimeStartPost)
		if err != nil {
			log.Fatal("Could not parse date", err)
			ctx.String(http.StatusBadRequest, "Could not parse date", gin.H{})
			return
		}
		dateTimeEnd, err := time.Parse("2006-01-02T15:04", dateTimeEindPost)
		if err != nil {
			log.Fatal("Could not parse date", err)
			ctx.String(http.StatusBadRequest, "Could not parse date", gin.H{})
			return
		}

		activity, err := database.GetActivityByUUID(uuid)
		if err != nil {
			fmt.Println("not found")
			ctx.HTML(http.StatusBadRequest, "", gin.H{})
		}
		activity.Message = message
		activity.Title = title
		activity.Location = location
		activity.DateTimeStart = dateTimeStart
		activity.DateTimeEnd = dateTimeEnd
		database.GormDB.Save(&activity)
		// Uses list here but works in this case
		activityList := []db.Activity{*activity}
		ctx.HTML(http.StatusOK, "activityListTemplate.html", gin.H{"activityList": activityList})
	})

	activityApi.POST("/submit", func(ctx *gin.Context) {
		title := ctx.PostForm("titel")
		message := ctx.PostForm("bericht")
		location := ctx.PostForm("locatie")
		dateTimeStartPost := ctx.PostForm("datumTijdStart")
		dateTimeEindPost := ctx.PostForm("datumTijdEnd")
		dateTimeStart, err := time.Parse("2006-01-02T15:04", dateTimeStartPost)
		if err != nil {
			log.Fatal("Could not parse start date: ", err)
			ctx.String(http.StatusBadRequest, "Could not parse date", gin.H{})
			return
		}
		dateTimeEnd, err := time.Parse("2006-01-02T15:04", dateTimeEindPost)
		if err != nil {
			log.Fatal("Could not parse end date: ", err)
			ctx.String(http.StatusBadRequest, "Could not parse date", gin.H{})
			return
		}
		newActivity := db.Activity{
			UUID:          uuid.New().String(),
			Message:       message,
			Title:         title,
			Location:      location,
			DateTimeStart: dateTimeStart,
			DateTimeEnd:   dateTimeEnd,
		}
		database.AddActivity(&newActivity)
		activityList = database.GetActivities()
		ctx.HTML(http.StatusOK, "activityListTemplateAdmin.html", gin.H{
			"activityList": activityList,
		})
	})
	activityApi.POST("/remove", func(ctx *gin.Context) {
		uuid := ctx.PostForm("UUID")
		database.RemoveActivityByUUID(uuid)
		activityList = database.GetActivities()
	})
	router.Run(":" + os.Getenv("PORT"))
}
