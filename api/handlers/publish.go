package handlers

import (
	"database/sql"

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

// Configures the slug (name of the site)
func SlugPublish(c *gin.Context) {
	var req struct {
		ID   string `json:"id" binding:"required"`
		Slug string `json:"slug" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}
	err := publish.SetSlug(req.Slug, req.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, req)
}

func SlugsPublish(c *gin.Context) {
	var req struct {
		Token string   `json:"token" binding:"required"`
		IDs   []string `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}
	_, err := utilities.GetJwtEmail(req.Token)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid token",
		})
		return
	}
	if len(req.IDs) == 0 {
		c.JSON(200, gin.H{})
	}
	siteSlugs := make(map[string]string)
	for _, id := range req.IDs {
		slug, err := publish.GetSiteSlug(id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		siteSlugs[id] = slug
	}
	c.JSON(200, siteSlugs)
}

func SitePublish(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
		Slug  string `json:"slug" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}
	_, err := utilities.GetJwtEmail(req.Token)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid token",
		})
		return
	}
	site, err := publish.GetSlug(req.Slug)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(200, gin.H{
				"code":    "NOTFOUND",
				"message": "Slug not found",
			})
			return
		}
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, site)
}
