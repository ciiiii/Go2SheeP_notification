package pusher

import (
	"github.com/ciiiii/Go2SheeP_notification/config"
	"github.com/gin-gonic/gin"
	notification "github.com/pusher/push-notifications-go"
)

var (
	client notification.PushNotifications
)

func init() {
	beamsClient, err := notification.New(config.Parser().Pusher.InstanceId, config.Parser().Pusher.PrivateKey)
	if err != nil {
		panic(err)
	}
	client = beamsClient

}

type NotifyRequest struct {
	Interests []string `json:"interests"`
	Icon      string   `json:"icon"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
}

func (n *NotifyRequest) send() (string, error) {
	if n.Icon == "" {
		n.Icon = "https://gotification.herokuapp.com/favicon.ico"
	}
	message := map[string]interface{}{
		"web": map[string]interface{}{
			"notification": map[string]interface{}{
				"icon":  n.Icon,
				"title": n.Title,
				"body":  n.Body,
			},
		},
	}
	pubId, err := client.PublishToInterests(n.Interests, message)
	if err != nil {
		return "", err
	}
	return pubId, nil
}

func NotifyHandler(c *gin.Context) {
	var notify NotifyRequest
	if err := c.ShouldBindJSON(&notify); err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	pubId, err := notify.send()
	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"pubId":   pubId,
	})
}
