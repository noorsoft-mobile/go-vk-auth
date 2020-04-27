package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-vk-api/vk"
	"golang.org/x/oauth2"
	vkAuth "golang.org/x/oauth2/vk"
	"log"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       []string{},
		Endpoint:     vkAuth.Endpoint,
	}
	r.GET("/", func(c *gin.Context) {
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		// получаем URL для редиректа на OAuth API VK и передаем его в темплейт
		c.HTML(http.StatusOK, "index.html", gin.H{
			"authUrl": url,
		})
	})
	r.GET("/auth", func(c *gin.Context) {
		ctx := context.Background()
		// получаем код от API VK из квери стринга
		authCode := c.Request.URL.Query()["code"]
		// меняем код на access токен
		tok, err := conf.Exchange(ctx, authCode[0])
		if err != nil {
			log.Fatal(err)
		}
		// создаем клиент для получения данных из API VK
		client, err := vk.NewClientWithOptions(vk.WithToken(tok.AccessToken))
		if err != nil {
			log.Fatal(err)
		}
		user := getCurrentUser(client)

		c.HTML(http.StatusOK, "auth.html", gin.H{
			"user": user,
		})
	})
	r.Run()
}

func getCurrentUser(api *vk.Client) User {
	var users []User

	_ = api.CallMethod("users.get", vk.RequestParams{
		"fields": "photo_400_orig,city",
	}, &users)
	return users[0]
}
