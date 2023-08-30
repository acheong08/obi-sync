package main

import (
	"github.com/acheong08/endless"
	"github.com/acheong08/obsidian-sync/api/handlers"
	"github.com/acheong08/obsidian-sync/database"
	gin "github.com/gin-gonic/gin"
)

func main() {
	dbConnection := database.NewDatabase()
	defer dbConnection.Close()
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Header("access-control-allow-origin", "*")
		c.Header("access-control-allow-methods", "GET, POST, OPTIONS")
		c.Header("access-control-allow-credentials", "true")
		// Allow all headers + content-type
		c.Header("access-control-allow-headers", "*, content-type, x-request-id")
		// c.Header("access-control-allow-headers", "content-type")

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
	userGroup.POST("signout", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})
	userGroup.POST("info", handlers.UserInfo)

	subscriptionGroup := router.Group("/subscription")
	subscriptionGroup.POST("list", handlers.ListSubscriptions)

	vaultGroup := router.Group("/vault")
	vaultGroup.POST("list", handlers.ListVaults)
	vaultGroup.POST("create", handlers.CreateVault)
	vaultGroup.POST("delete", handlers.DeleteVault)
	vaultGroup.POST("access", handlers.AccessVault)

	vaultShareGroup := router.Group("/vault/share")
	vaultShareGroup.POST("list", handlers.ListVaultShares)
	vaultShareGroup.POST("invite", handlers.InviteVaultShare)
	vaultShareGroup.POST("remove", handlers.RemoveVaultShare)

	publishAPI := router.Group("/api")
	publishGroup := router.Group("/publish")
	// Checks if a site already exists
	publishAPI.POST("site", handlers.SitePublish)
	// Creates a new site (on host. similar to vault)
	publishGroup.POST("create", handlers.CreatePublish)
	// list files in a site
	publishGroup.POST("list", handlers.ListPublish)
	publishAPI.POST("list", handlers.ListPublish)
	// Configures the slug (name of the site)
	publishAPI.POST("slug", handlers.SlugPublish)
	// returns a map of slug id to name
	publishAPI.POST("slugs", handlers.SlugsPublish)
	publishAPI.POST("upload", handlers.UploadFile)
	publishAPI.POST("remove", handlers.RemoveFile)

	router.GET("/published/:slug/*path", handlers.GetPublishedFile)

	router.GET("/", handlers.WsHandler)
	router.GET("/ws", handlers.WsHandler)

	endless.ListenAndServe("127.0.0.1:3000", router)

}
