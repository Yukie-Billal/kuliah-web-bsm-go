package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"success": true,
		"data":    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(201, gin.H{
		"success": true,
		"data":    data,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(404, gin.H{
		"success": false,
		"message": message,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(400, gin.H{
		"success": false,
		"message": message,
	})
}

func InternalServerError(c *gin.Context, message ...string) {
	msg := "internal server error"

	if len(message) > 0 {
		msg = message[0]
	}
	c.JSON(500, gin.H{
		"success": false,
		"message": msg,
	})
}
