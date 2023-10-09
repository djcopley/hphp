package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
)

var (
	redisDB *redis.Client
	ctx     context.Context
)

var scoreResponseString = `<p id="%s-score" class="text-center">Score: %d</p>`

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
	stringScore, err := redisDB.Get(ctx, houseName).Result()
	if err != nil {
		panic(err)
	}
	score, err := strconv.Atoi(stringScore)
	if err != nil {
		panic(err)
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, fmt.Sprintf(scoreResponseString, houseName, score))
}

func setScore(c *gin.Context) {
	houseName := c.Param("houseName")
	newScore, err := strconv.Atoi(c.PostForm("newScore"))
	if err != nil {
		panic(err.Error())
	}
	if err = redisDB.Set(ctx, houseName, newScore, 0).Err(); err != nil {
		panic(err.Error())
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, fmt.Sprintf(scoreResponseString, houseName, newScore))
}

func patchScore(c *gin.Context) {
	houseName := c.Param("houseName")
	incrementBy, err := strconv.Atoi(c.PostForm("incrementBy"))
	if err != nil {
		panic(err.Error())
	}
	newValue, err := redisDB.IncrBy(ctx, houseName, int64(incrementBy)).Result()
	if err != nil {
		panic(err.Error())
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, fmt.Sprintf(scoreResponseString, houseName, newValue))
}

func main() {
	initRedis("localhost:6379", "")
	r := gin.Default()

	r.GET("/score/:houseName", getScore)
	r.POST("/score/:houseName", setScore)
	r.PATCH("/score/:houseName", patchScore)

	// Serve static files from the "web" directory.
	r.StaticFS("/static", http.Dir("web/static"))
	// Define a catch-all route to serve the index.html file.
	r.NoRoute(func(c *gin.Context) {
		c.File("web/index.html")
	})

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
