package config

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

func InitDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		panic("failed to load .env file")
	}

	dsn := os.Getenv("DATABASE_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	return db
}

func InitRedis() *redis.Client {
	err := godotenv.Load()
	if err != nil {
		panic("failed to load .env file")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		panic("failed to connect Redis: " + err.Error())
	}

	return client
}
