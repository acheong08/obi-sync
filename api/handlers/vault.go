package handlers

import (
	"github.com/acheong08/obsidian-sync/database/vault"
	"github.com/acheong08/obsidian-sync/utilities"
	"github.com/gin-gonic/gin"
	password_generator "github.com/sethvargo/go-password/password"
)

func ListVaults(c *gin.Context) {
	type request struct {
		Token string `json:"token" binding:"required"`
	}
	type response struct {
		Shared []*vault.Vault `json:"shared"`
		Vaults []*vault.Vault `json:"vaults"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	email, err := utilities.GetJwtEmail(req.Token)
	if err != nil {
		// Unauthorized
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	vaults, err := vault.GetVaults(email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Get shared vaults
	shared, err := vault.GetSharedVaults(email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response{
		Shared: shared,
		Vaults: vaults,
	})

}

func CreateVault(c *gin.Context) {
	type request struct {
		KeyHash string `json:"keyhash"`
		Name    string `json:"name" binding:"required"`
		Salt    string `json:"salt"`
		Token   string `json:"token" binding:"required"`
	}
	// Response is vault details
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	email, err := utilities.GetJwtEmail(req.Token)
	if err != nil {
		// Unauthorized
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	var password string
	var salt string
	var keyHash string
	// Generate password if keyhash is not provided
	if req.Salt == "" {

		password, err = password_generator.Generate(20, 5, 5, false, true)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		salt, err = password_generator.Generate(20, 5, 5, false, true)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		keyHash = ""
	} else {
		salt = req.Salt
		if req.KeyHash != "" {
			keyHash = req.KeyHash
		} else {
			c.JSON(400, gin.H{"error": "keyhash must be provided if salt is provided"})
		}
	}
	vault, err := vault.NewVault(req.Name, email, password, salt, keyHash)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, vault)

}

func DeleteVault(c *gin.Context) {
	type request struct {
		Token    string `json:"token" binding:"required"`
		VaultUID string `json:"vault_uid" binding:"required"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	email, err := utilities.GetJwtEmail(req.Token)
	if err != nil {
		// Unauthorized
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	err = vault.DeleteVault(req.VaultUID, email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{})
}

func AccessVault(c *gin.Context) {
	type request struct {
		Host     string `json:"host" binding:"required"`
		KeyHash  string `json:"keyhash" binding:"required"`
		Token    string `json:"token" binding:"required"`
		VaultUID string `json:"vault_uid" binding:"required"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	email, err := utilities.GetJwtEmail(req.Token)
	if err != nil {
		// Unauthorized
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	if !vault.HasAccessToVault(req.VaultUID, email) {
		c.JSON(401, gin.H{"error": "You do not have access to this vault"})
		return
	}
	_, err = vault.GetVault(req.VaultUID, req.KeyHash)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}
	// Get user details
	user, err := vault.UserInfo(email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"allowed": true,
		"email":   email,
		"name":    user.Name,
		"useruid": "b094fc51bf40b9ddb9ff43d4aadfa962", // Not necessary...
	})

}
