package utilities

import (
	"errors"

	"github.com/acheong08/obsidian-sync/config"
	"github.com/golang-jwt/jwt/v5"
)

func GetJwtEmail(jwtString string) (string, error) {
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
