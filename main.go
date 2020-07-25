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
