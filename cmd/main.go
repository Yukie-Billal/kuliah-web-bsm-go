package main

import (
	"kuliah-web-bsm-go/internal/config"
	"kuliah-web-bsm-go/internal/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadEnv()

	db := config.ConnectDB()

	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1", "localhost"})

	routes.Register(router, db)

	router.Run(":" + os.Getenv("APP_PORT"))
}
