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
		c.Header("access-control-allow-methods", "GET, POST")

		// Add database connection to context
		c.Set("db", dbConnection)
	})
	userGroup := router.Group("/user")
	userGroup.POST("signin", handlers.Signin)
}
