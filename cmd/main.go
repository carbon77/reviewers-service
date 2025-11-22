package main

import (
	"fmt"
	"reviewers/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello",
		})
	})

	addr := fmt.Sprintf(":%s", cfg.HttpPort)
	router.Run(addr)
}
