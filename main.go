// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"project-root/api"
	"project-root/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func SubscribeToNotifications(redisClient *redis.Client) {
	ctx := context.Background()
	pubsub := redisClient.Subscribe(ctx, "notifications")
	defer pubsub.Close()

	fmt.Println("Subscribed to notifications...") // Добавьте это сообщение

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			fmt.Printf("Error receiving message: %v\n", err)
			continue
		}
		fmt.Printf("Received notification: %s\n", msg.Payload)
	}
}

func main() {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	if err := config.InitDB(); err != nil {
		panic(fmt.Sprintf("Failed to connect to PostgreSQL: %v", err))
	}
	config.InitRedis()

	router := gin.Default()
	// Настройка CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:8081"}, // Разрешить запросы с вашего фронтенда
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	api.RegisterBookingRoutes(router)
	api.RegisterBrigadeRoutes(router)
	api.RegisterNotificationRoutes(router)
	go SubscribeToNotifications(config.RedisClient)

	go func() {
		pubsub := config.RedisClient.Subscribe(context.Background(), "booking_notifications")
		defer pubsub.Close()

		for msg := range pubsub.Channel() {
			fmt.Printf("Received notification: %s\n", msg.Payload)
		}
	}()

	router.Run(":8080")
}
