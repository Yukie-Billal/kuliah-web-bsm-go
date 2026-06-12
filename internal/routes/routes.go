package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(
	router *gin.Engine,
	db *gorm.DB,
) {

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
