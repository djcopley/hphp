package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	redisDB *redis.Client
	ctx     context.Context
)

func initRedis(addr string, password string) {
	redisDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx = context.Background()

	pong, err := redisDB.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis:", pong)
}

func getScore(c *gin.Context) {
	houseName := c.Param("houseName")
	value, err := redisDB.Get(ctx, houseName).Result()
	if err != nil {
		panic(err)
	}
	scoreString := fmt.Sprintf("<p class=\"text-center\">Score: %s</p>", value)
	c.String(http.StatusOK, scoreString)
}

func setScore(c *gin.Context) {
	houseName := c.Param("houseName")
	var newScore struct {
		Score int `json:"score"`
	}
	if err := c.ShouldBindJSON(&newScore); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	fmt.Println(newScore)
	err := redisDB.Set(ctx, houseName, newScore.Score, 0).Err()
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func main() {
	initRedis("localhost:6379", "")
	r := gin.Default()

	r.GET("/score/:houseName", getScore)
	r.POST("/score/:houseName", setScore)

	// Serve static files from the "website" directory.
	r.StaticFS("/static", http.Dir("website/static"))
	// Define a catch-all route to serve the index.html file.
	r.NoRoute(func(c *gin.Context) {
		c.File("website/index.html")
	})

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
