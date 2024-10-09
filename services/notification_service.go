package services

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Notification struct {
	Message string `json:"message"`
}

// Функция для отправки уведомления
func SendNotification(redisClient *redis.Client, notification Notification) error {
	ctx := context.Background()
	err := redisClient.Publish(ctx, "notifications", notification.Message).Err()
	if err != nil {
		return fmt.Errorf("could not publish notification: %v", err)
	}
	return nil
}
