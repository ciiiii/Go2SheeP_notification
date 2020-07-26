package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ciiiii/Go2SheeP_notification/config"
	"github.com/ciiiii/Go2SheeP_notification/pusher"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var conf *oauth2.Config

type UserInfo struct {
	Email   string `json:"email"`
	Picture string `json:"picture"`
	Sub     string `json:"sub"`
}

func init() {
	conf = &oauth2.Config{
		ClientID:     config.Parser().OAuth.ClientId,
		ClientSecret: config.Parser().OAuth.ClientSecret,
		RedirectURL:  config.Parser().OAuth.RedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}

func randState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func authHandler(c *gin.Context) {
	retrievedState, err := c.Cookie("state")
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("get state faild"))
		return
	}
	if retrievedState != c.Query("state") {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid state: %s", retrievedState))
		return
	}

	token, err := conf.Exchange(context.Background(), c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	client := conf.Client(context.Background(), token)
	email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	defer email.Body.Close()
	data, _ := ioutil.ReadAll(email.Body)
	var userInfo UserInfo
	json.Unmarshal(data, &userInfo)
	c.SetCookie("picture", userInfo.Picture, 3600, "/", "localhost", false, false)
	c.SetCookie("email", userInfo.Email, 3600, "/", "localhost", false, false)
	c.SetCookie("token", token.AccessToken, 3600, "/", "localhost", false, false)
	c.Redirect(301, "/")
}

func configHandler(c *gin.Context) {
	state := randState()
	c.SetCookie("state", state, 3600, "/", "localhost", false, false)
	c.JSON(http.StatusOK, gin.H{
		"url": conf.AuthCodeURL(state),
	})
}

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.DisableConsoleColor()

	r.LoadHTMLFiles("static/index.html")
	r.Static("/css", "./static/css")
	r.Static("/js", "./static/js")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.StaticFile("/service-worker.js", "./static/service-worker.js")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/auth", authHandler)
	r.GET("/config", configHandler)

	r.GET("/notify", pusher.NotifyHandler)
	server := &http.Server{
		Addr:    ":" + config.Parser().App.Port,
		Handler: r,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
