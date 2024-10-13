// api/notification.go
package api

import (
	"context"
	"io"
	"net/http"
	"project-root/config"

	"github.com/gin-gonic/gin"
)

func RegisterNotificationRoutes(router *gin.Engine) {
	router.GET("/api/notifications", getNotifications)
	router.GET("/api/notifications/stream", streamNotifications) // Новый маршрут для потока уведомлений
}

func streamNotifications(c *gin.Context) {
	ctx := context.Background()
	pubsub := config.RedisClient.Subscribe(ctx, "notifications")
	defer pubsub.Close()

	// Устанавливаем заголовок Content-Type для SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// Начинаем поток
	c.Stream(func(w io.Writer) bool {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			return false
		}
		// Отправляем новое уведомление клиенту
		c.SSEvent("message", msg.Payload)
		return true
	})
}
func getNotifications(c *gin.Context) {
	notifications, err := config.RedisClient.LRange(context.Background(), "notifications", 0, -1).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve notifications"})
		return
	}
	c.JSON(http.StatusOK, notifications) // Вернем массив уведомлений
}
