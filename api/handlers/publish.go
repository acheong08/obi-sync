package handlers

import (
	"net/url"

	"github.com/acheong08/obsidian-sync/database/publish"
	"github.com/acheong08/obsidian-sync/utilities"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListPublish(c *gin.Context) {
	var req struct {
		Token   string `json:"token" binding:"required"`
		Version int    `json:"version"`
		ID      string `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
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
		return
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

func CreatePublish(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
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
	// Check how many sites the user has
	sites, err := publish.GetSites(email)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	if len(sites) >= 1 {
		c.JSON(200, gin.H{
			"error": "You have reached the limit of 1 site",
		})
		return
	}
	site, err := publish.CreateSite(email)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, site)
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
		if err == gorm.ErrRecordNotFound {
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

func RemoveFile(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
		Path  string `json:"path" binding:"required"`
		ID    string `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}
	email, err := utilities.GetJwtEmail(req.Token)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "invalid token",
		})
		return
	}
	siteOwner, err := publish.GetSiteOwner(req.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	if siteOwner != email {
		c.JSON(403, gin.H{
			"error": "You do not have permission to delete this file",
		})
		return
	}
	err = publish.RemoveFile(req.ID, req.Path)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{})
}

func UploadFile(c *gin.Context) {
	Token := c.Request.Header.Get("obs-token")
	email, err := utilities.GetJwtEmail(Token)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid token",
		})
		return
	}
	var file publish.File = publish.File{
		Size: c.Request.ContentLength,
		Hash: c.Request.Header.Get("obs-hash"),
		Slug: c.Request.Header.Get("obs-id"),
		Path: c.Request.Header.Get("obs-path"),
	}
	// Path is URL encoded. Unencode it
	file.Path, err = url.QueryUnescape(file.Path)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	siteOwner, err := publish.GetSiteOwner(file.Slug)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	if siteOwner != email {
		c.JSON(403, gin.H{
			"error": "You do not have permission to upload to this site",
		})
	}
	// Read body as text
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	file.Data = data
	err = publish.NewFile(&file)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{})

}

func GetPublishedFile(c *gin.Context) {
	// Get slug and path from url
	slug := c.Param("slug")
	path := c.Param("path")[1:]

	// Get site id from slug
	site, err := publish.GetSlug(slug)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"error": "Site not found",
			})
			return
		}
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Get file from site id and path
	file, err := publish.GetFile(site.ID, path)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
			"site":  site,
			"path":  path,
		})
		return
	}
	// Return file []byte
	c.Data(200, "text/markdown; charset=utf-8", file)
}
