package handlers

import (
	"github.com/acheong08/obsidian-sync/config"
	"github.com/acheong08/obsidian-sync/database"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

func Signin(c *gin.Context) {
	c.Header("access-control-allow-credentials", "true")
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
	dbConnection := c.MustGet("db").(*database.Database)
	userInfo, err := dbConnection.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": userInfo.Email,
	}).SignedString(config.Secret)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, response{
		Email:   userInfo.Email,
		License: userInfo.License,
		Name:    userInfo.Name,
		Token:   token,
	})

}
