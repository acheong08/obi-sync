package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

func B(e []byte) []byte {
	return e[:]
}

func j(e []byte) ([]byte, error) {
	hash := sha256.Sum256(e)
	return hash[:], nil
}

func G(e []byte) string {
	t := make([]string, len(e))
	for _, r := range e {
		t = append(t, fmt.Sprintf("%x", r>>4))
		t = append(t, fmt.Sprintf("%x", 15&r))
	}
	return strings.Join(t, "")
}

func Y(e []byte) (cipher.Block, error) {
	return aes.NewCipher(e)
}

const (
	K    = "aes-256-gcm"
	cost = 32768
)

type scrypt_Params struct {
	N       int
	R       int
	P       int
	KeyLen  int
	SaltLen int
	MaxMem  int
	MaxSalt int
}

var Q = scrypt_Params{
	N:       cost,
	R:       8,
	P:       1,
	MaxMem:  67108864,
	MaxSalt: 32,
}

func J(e, t []byte) ([]byte, error) {
	normalizedE := normalizeNFKC(string(e))
	normalizedT := normalizeNFKC(string(t))
	key, err := scrypt.Key([]byte(normalizedE), []byte(normalizedT), Q.N, Q.R, Q.P, Q.KeyLen)
	if err != nil {
		return nil, err
	}
	return B(key), nil
}

func makeKeyHash(e, t []byte) (string, error) {
	n, err := J(e, t)
	if err != nil {
		return "", err
	}
	hash, err := j(n)
	if err != nil {
		return "", err
	}
	return G(hash), nil
}

func normalizeNFKC(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func main() {
	e := []byte("ZsSjgKx4yaeBNCFipS)T")
	t := []byte("jePEuEPhNsr8zguY3%98")
	hash, err := makeKeyHash(e, t)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(hash)
}
