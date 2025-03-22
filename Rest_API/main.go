package main

import (
	"github.com/gin-gonic/gin"
	"rest-api/database"
	"rest-api/routes"
)

func main() {
	database.InitDB()
	// Initialize HTTP server with basic features
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run("localhost:8080")
}
