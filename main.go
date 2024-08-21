package main

import (
	"myapp/config"
	"myapp/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize MySQL and Redis
	config.ConnectDatabase()
	config.ConnectRedis()

	// Initialize the custom logger
	config.InitLogger()

	// Set up the Gin router
	router := gin.Default()

	// Initialize routes
	routes.InitializeRoutes(router)

	// Start the server
	router.Run(":8080")
}
