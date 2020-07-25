package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.DisableConsoleColor()

	r.LoadHTMLGlob("static/*.html")
	r.Static("/css", "./static/css")
	r.Static("/js", "./static/js")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.StaticFile("/service-worker.js", "./static/service-worker.js")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
