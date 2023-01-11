package main

import (
	"context"
	"net/http"
	"sum/pkg/infrastructure"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := infrastructure.NewLogger()
	gormLogger := infrastructure.NewGormLogger(logger)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	gormLogger.Info(context.Background())
	r.Run(":8080")

}
