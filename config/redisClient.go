package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

var Redis *redis.Client
var ctx = context.Background()

func ConnectRedis() *redis.Client {
	address := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Connection fail in Redis：", pong, err)
		return nil
	}
	fmt.Println("Connection success in Redis：", pong)
	Redis = client
	return client
}
