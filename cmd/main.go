package main

import (
	"fmt"
	"reviewers/internal/config"
	"reviewers/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	conn := db.Connect(cfg)

	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello",
		})
	})

	router.GET("/users", func(c *gin.Context) {
		var users []db.User
		conn.Find(&users)
		c.JSON(200, users)
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	router.Run(addr)
}
