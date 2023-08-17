package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/scrypt"
)

func B(e []byte) []byte {
	return e
}

func j(e []byte) []byte {
	digest := sha256.Sum256(e)
	return digest[:]
}

const (
	cost = 32768
	r    = 8
	p    = 1
)

func J(e, t string) ([]byte, error) {
	normalizedE := []byte(e)
	normalizedT := []byte(t)

	key, err := scrypt.Key(normalizedE, normalizedT, cost, r, p, 32)
	if err != nil {
		return nil, err
	}

	return B(key), nil
}

func makeKeyHash(e, t string) (string, error) {
	n, err := J(e, t)
	if err != nil {
		return "", err
	}

	hash := j(n)
	fmt.Println(hash)
	return hex.EncodeToString(hash), nil
}

func main() {
	hash, err := makeKeyHash("ZsSjgKx4yaeBNCFipS)T", "jePEuEPhNsr8zguY3%98")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(hash)
}
