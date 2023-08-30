package handlers

import (
	"strings"
	"time"

	"github.com/acheong08/obsidian-sync/config"
	"github.com/acheong08/obsidian-sync/database/vault"
	"github.com/acheong08/obsidian-sync/utilities"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func Signin(c *gin.Context) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		Email   string `json:"email"`
		License string `json:"license"`
		Name    string `json:"name"`
		Token   string `json:"token"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userInfo, err := vault.Login(req.Email, req.Password)
	if err != nil {
		// 200 because the app doesn't check the status code.
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": userInfo.Email,
	}).SignedString(config.Secret)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response{
		Email:   userInfo.Email,
		License: userInfo.License,
		Name:    userInfo.Name,
		Token:   token,
	})

}

func UserInfo(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"error": "not logged in"})
		return
	}
	email, err := utilities.GetJwtEmail(req.Token)
	if err != nil {
		c.JSON(200, gin.H{"error": "not logged in"})
		return
	}
	userInfo, err := vault.UserInfo(email)
	if err != nil {
		c.JSON(200, gin.H{"error": "not logged in"})
		return
	}
	c.JSON(200, gin.H{
		"uid":     uuid.New().String(),
		"email":   email,
		"name":    userInfo.Name,
		"payment": "",
		"license": "",
		"credit":  0,
		"mfa":     false,
		"discount": gin.H{
			"status":    "approved",
			"expiry_ts": time.Now().UnixMilli() + time.Hour.Milliseconds()*24*365,
			"type":      "education",
		},
	})
}

func SignUp(c *gin.Context) {
	var req struct {
		Email     string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required"`
		FullName  string `json:"name" binding:"required"`
		SignUpKey string `json:"signup_key"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if req.SignUpKey != config.SignUpKey && config.SignUpKey != "" {
		c.JSON(400, gin.H{"error": "invalid signup key"})
		return
	}

	err := vault.NewUser(req.Email, strings.TrimSpace(req.Password), req.FullName)

	if err != nil {
		c.JSON(500, gin.H{"error": "not sign up!"})
		return
	}

	c.JSON(200, gin.H{
		"email": req.Email,
		"name":  req.FullName,
	})
}
