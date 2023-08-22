package main

import (
	"github.com/acheong08/obsidian-sync/api/handlers"
	"github.com/acheong08/obsidian-sync/database"
	gin "github.com/gin-gonic/gin"
)

func main() {
	dbConnection := database.NewDatabase()
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Header("access-control-allow-origin", "app://obsidian.md")
		c.Header("access-control-allow-methods", "GET, POST, OPTIONS")
		c.Header("access-control-allow-credentials", "true")
		// Allow all headers
		c.Header("access-control-allow-headers", "*")

		// Add database connection to context
		c.Set("db", dbConnection)
	})
	// Respond to all OPTIONS requests with 200
	router.OPTIONS("/*cors", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	userGroup := router.Group("/user")
	userGroup.POST("signin", handlers.Signin)

	subscriptionGroup := router.Group("/subscription")
	subscriptionGroup.POST("list", handlers.ListSubscriptions)

	vaultGroup := router.Group("/vault")
	vaultGroup.POST("list", handlers.ListVaults)
	vaultGroup.POST("create", handlers.CreateVault)
	vaultGroup.POST("delete", handlers.DeleteVault)

	router.Run(":3000")
}
