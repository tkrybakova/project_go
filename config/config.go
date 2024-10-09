// config/config.go
package config

import (
    "context"
    "fmt"
    "os"

    "github.com/jackc/pgx/v4"
    "github.com/go-redis/redis/v8"
)

var DB *pgx.Conn
var RedisClient *redis.Client

func InitDB() error {
    dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASS"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )

    var err error
    DB, err = pgx.Connect(context.Background(), dsn)
    if err != nil {
        return fmt.Errorf("unable to connect to database: %v", err)
    }

    fmt.Println("Connected to PostgreSQL!")
    return nil
}

func InitRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
        Password: os.Getenv("REDIS_PASS"),
        DB:       0,
    })
    fmt.Println("Connected to Redis!")
}
