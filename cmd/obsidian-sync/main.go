package main

import (
	"github.com/acheong08/endless"
	"github.com/acheong08/obsidian-sync/api/handlers"
	"github.com/acheong08/obsidian-sync/config"
	gin "github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Header("access-control-allow-origin", "*")
		c.Header("access-control-allow-methods", "GET, POST, OPTIONS")
		c.Header("access-control-allow-credentials", "true")
		// Allow all headers + content-type
		c.Header("access-control-allow-headers", "*, content-type, x-request-id")
	})
	// Respond to all OPTIONS requests with 200
	router.OPTIONS("/*cors", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	userGroup := router.Group("/user")
	userGroup.POST("signup", handlers.SignUp)
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
	publishAPI.POST("site", handlers.SiteInfo)
	publishGroup.POST("create", handlers.CreateSite)
	publishGroup.POST("list", handlers.ListSites)
	publishAPI.POST("list", handlers.ListSites)
	publishAPI.POST("slug", handlers.ConfigureSiteSlug)
	publishAPI.POST("slugs", handlers.GetSlugInfo)
	publishAPI.POST("upload", handlers.UploadFile)
	publishAPI.POST("remove", handlers.RemoveFile)
	publishGroup.POST("delete", handlers.DeleteSite)

	router.GET("/published/:slug/*path", handlers.GetPublishedFile)
	router.GET("/published/:slug", handlers.GetSiteIndex)

	router.GET("/", handlers.WsHandler)
	router.GET("/ws", handlers.WsHandler)
	router.GET("/ws.obsidian.md", handlers.WsHandler)

	endless.ListenAndServe(config.AddressHttp, router)

}
