package db

import (
	"chat-system/internal/utils"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		fmt.Println("Invalid Redis DB: %v", err)
		utils.GetLogger().Error("Invalid Redis DB: ", err)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	_, err = RedisClient.Ping(Ctx).Result()
	if err != nil {
		utils.GetLogger().Error("Failed to connect to Redis: ", err)
		fmt.Println("Failed to connect to Redis: %v", err)
	}
}
