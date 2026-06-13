package routes

import (
	"kuliah-web-bsm-go/internal/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello world",
	})
}

func Register(
	router *gin.Engine,
	db *gorm.DB,
) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	router.GET("/users", userHandler.Index)
	router.GET("/users/:id", userHandler.Show)
	router.POST("/users", userHandler.Store)
	router.PUT("/users/:id", userHandler.Update)
	router.DELETE("/users/:id", userHandler.Delete)

	router.GET("/hello", Test)
}
