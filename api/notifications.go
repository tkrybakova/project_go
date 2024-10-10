// api/notification.go
package api

import (
	"context"
	"net/http"
	"project-root/config"

	"github.com/gin-gonic/gin"
)

func RegisterNotificationRoutes(router *gin.Engine) {
	router.GET("/api/notifications", getNotifications)
}

func getNotifications(c *gin.Context) {
	notifications, err := config.RedisClient.LRange(context.Background(), "notifications", 0, -1).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve notifications"})
		return
	}
	c.JSON(http.StatusOK, notifications) // Вернем массив уведомлений
}
