package handlers

import (
	"github.com/acheong08/obsidian-sync/publish"
	"github.com/acheong08/obsidian-sync/utilities"
	"github.com/gin-gonic/gin"
)

func ListPublish(c *gin.Context) {
	var req struct {
		Token   string `json:"token" binding:"required"`
		Version string `json:"version"`
		ID      string `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}
	email, err := utilities.GetJwtEmail(req.Token)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid token",
		})
	}
	if req.ID == "" {
		sites, err := publish.GetSites(email)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"sites":  sites,
			"shared": make([]interface{}, 0),
			"limit":  1,
		})
		return
	}
	siteOwner, err := publish.GetSiteOwner(req.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
	files, err := publish.GetFiles(req.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"files": files,
		"owner": siteOwner == email,
	})

}
