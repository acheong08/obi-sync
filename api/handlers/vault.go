package handlers

import (
	"errors"

	"github.com/acheong08/obsidian-sync/config"
	"github.com/acheong08/obsidian-sync/database"
	"github.com/acheong08/obsidian-sync/vault"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

func ListVaults(c *gin.Context) {
	type request struct {
		Token string `json:"token"`
	}
	type response struct {
		Shared []any          `json:"shared"`
		Vaults []*vault.Vault `json:"vaults"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	email, err := getJwtEmail(req.Token)
	if err != nil {
		// Unauthorized
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	dbConnection := c.MustGet("db").(*database.Database)
	vaults, err := dbConnection.GetVaults(email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response{
		Shared: []any{},
		Vaults: vaults,
	})

}

func getJwtEmail(jwtString string) (string, error) {
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		return config.Secret, nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token")
	}
	email, ok := claims["email"].(string)
	if !ok {
		return "", errors.New("invalid token")
	}
	return email, nil
}
