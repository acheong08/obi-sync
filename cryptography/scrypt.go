package cryptography

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/scrypt"
)

const (
	cost = 32768
	r    = 8
	p    = 1
)

func getKey(e, t string) ([]byte, error) {
	normalizedE := []byte(e)
	normalizedT := []byte(t)

	key, err := scrypt.Key(normalizedE, normalizedT, cost, r, p, 32)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func MakeKeyHash(e, t string) (string, error) {
	n, err := getKey(e, t)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(n)
	return hex.EncodeToString(hash[:]), nil
}
