package routes

import (
	"fmt"
	"github.com/djcopley/hphp/internal/db"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
)

var scoreResponseString = `<p id="%s-score" class="text-center">Score: %d</p>`

func GetScore(c *gin.Context) {
	houseName := c.Param("houseName")
	score := db.GetScore(houseName)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, fmt.Sprintf(scoreResponseString, houseName, score))
}

func SetScore(c *gin.Context) {
	houseName := c.Param("houseName")
	newScore, err := strconv.Atoi(c.PostForm("newScore"))
	if err != nil {
		panic(err.Error())
	}

	broadcastEvent(fmt.Sprintf("%d points for %s", newScore, houseName))
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, fmt.Sprintf(scoreResponseString, houseName, newScore))
}

func PatchScore(c *gin.Context) {
	houseName := c.Param("houseName")
	incrementBy, err := strconv.Atoi(c.PostForm("incrementBy"))
	if err != nil {
		log.Panicf("Failed to convert incrementBy to int: %v", err)
	}
	newScore := db.IncrementScore(houseName, incrementBy)
	broadcastEvent(fmt.Sprintf("%d points for %s", newScore, houseName))
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, fmt.Sprintf(scoreResponseString, houseName, newScore))
}

func GetHouseEvents(c *gin.Context) {
	sseClient := make(chan string)
	addSSEClient(sseClient)

	clientClosed := c.Writer.CloseNotify()
	go func() {
		<-clientClosed
		removeSSEClient(sseClient)
		close(sseClient)
	}()

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-sseClient; ok {
			c.SSEvent("scores_updated", msg)
			return true
		}
		return false
	})
}
