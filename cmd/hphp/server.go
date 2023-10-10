package main

import (
	"github.com/djcopley/hphp/internal/db"
	"github.com/djcopley/hphp/internal/routes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	db.InitRedis("localhost:6379", "")
	r := gin.Default()

	r.GET("/house-events", routes.GetHouseEvents)
	r.GET("/score/:houseName", routes.GetScore)
	r.POST("/score/:houseName", routes.SetScore)
	r.PATCH("/score/:houseName", routes.PatchScore)

	// Serve static files from the "web" directory.
	r.StaticFS("/static", http.Dir("web/static"))
	// Define a catch-all route to serve the index.html file.
	r.NoRoute(func(c *gin.Context) {
		c.File("web/index.html")
	})

	if err := r.Run(":8080"); err != nil {
		log.Panicf("Failed to start gin server: %v", err)
	}
}
