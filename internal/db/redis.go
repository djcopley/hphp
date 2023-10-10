package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
)

var (
	redisDB *redis.Client
	ctx     context.Context
)

func InitRedis(addr string, password string) {
	redisDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx = context.Background()

	pong, err := redisDB.Ping(ctx).Result()
	if err != nil {
		log.Panicf("Failed to connect to redis: %v", err.Error())
	}
	log.Println("Connected to Redis:", pong)
}

func GetScore(houseName string) int {
	stringScore, err := redisDB.Get(ctx, houseName).Result()
	if err != nil {
		log.Panicf("Failed to get score from redis: %v", err)
	}
	score, err := strconv.Atoi(stringScore)
	if err != nil {
		log.Panicf("Failed to convert house %s score to int: %v", houseName, err)
	}
	return score
}

func IncrementScore(houseName string, incrementBy int) int {
	newScore, err := redisDB.IncrBy(ctx, houseName, int64(incrementBy)).Result()
	if err != nil {
		log.Panicf("Failed to increment score for house %s: %v", houseName, err)
	}
	return int(newScore)
}
