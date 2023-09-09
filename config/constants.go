package config

import (
	"crypto/rand"
	"encoding/gob"
	"log"
	"os"
	"path"
	"strconv"

	dotenv "github.com/joho/godotenv"
)

var SecretPath = "secret.gob"

var Host = "localhost:3000"

var AddressHttp = "127.0.0.1:3000"

var DataDir = "."

var Secret []byte

var SignUpKey string

var MaxStorageBytes int64 = 10 * 1073741824

var MaxSitesPerUser int = 5

// Generate a random password, hash it, and store it in the Secret variable & a file
// Load secret.gob if it exists
func init() {
	err := dotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	if os.Getenv("HOST") != "" { // Legacy. Use DOMAIN_NAME instead
		log.Println("Warning: HOST is deprecated. Use DOMAIN_NAME instead")
		Host = os.Getenv("HOST")
	}
	if os.Getenv("DOMAIN_NAME") != "" {
		Host = os.Getenv("DOMAIN_NAME")
	}
	if os.Getenv("ADDR_HTTP") != "" {
		AddressHttp = os.Getenv("ADDR_HTTP")
	}
	if os.Getenv("SIGNUP_KEY") != "" {
		SignUpKey = os.Getenv("SIGNUP_KEY")
	}
	if os.Getenv("DATA_DIR") != "" {
		DataDir = os.Getenv("DATA_DIR")

		if _, err := os.Stat(DataDir); os.IsNotExist(err) {
			err := os.Mkdir(DataDir, os.ModePerm)

			if err != nil {
				panic(err)
			}
		}

		SecretPath = path.Join(DataDir, "secret.gob")
	}
	if os.Getenv("MAX_STORAGE_GB") != "" {
		gb_size, err := strconv.Atoi(os.Getenv("MAX_STORAGE_GB"))
		if err != nil {
			panic(err)
		}
		MaxStorageBytes = int64(gb_size) * 1073741824
	}
	if os.Getenv("MAX_SITES_PER_USER") != "" {
		sites, err := strconv.Atoi(os.Getenv("MAX_SITES_PER_USER"))
		if err != nil {
			panic(err)
		}
		MaxSitesPerUser = sites
	}
	if _, err := os.Stat(SecretPath); err != nil {
		Secret = make([]byte, 64)
		rand.Read(Secret)
		f, err := os.Create(SecretPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		encoder := gob.NewEncoder(f)
		err = encoder.Encode(Secret)
		if err != nil {
			panic(err)
		}
	} else {
		f, err := os.Open(SecretPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		decoder := gob.NewDecoder(f)
		err = decoder.Decode(&Secret)
		if err != nil {
			panic(err)
		}
	}

}
