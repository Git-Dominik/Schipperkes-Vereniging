package main

import (
	"net/http"
	"os"
	"slices"

	"github.com/gin-gonic/gin"
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

	router.Run(":8080")
}
