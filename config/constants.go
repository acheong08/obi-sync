package config

import (
	"crypto/rand"
	"encoding/gob"
	"os"
	"path"
)

var DBPath = "database.db"

var SecretPath = "secret.gob"

var Host = "localhost:3000"

var AddressHttp = "0.0.0.0:3000"

var DataDir = "."

var Secret []byte

// Generate a random password, hash it, and store it in the Secret variable & a file
// Load secret.gob if it exists
func init() {
	if os.Getenv("HOST") != "" {
		Host = os.Getenv("HOST")
	}
	if os.Getenv("ADDR_HTTP") != "" {
		AddressHttp = os.Getenv("ADDR_HTTP")
	}
	if os.Getenv("DATA_DIR") != "" {
		DataDir = os.Getenv("DATA_DIR")

		if _, err := os.Stat(DataDir); os.IsNotExist(err) {
			err := os.Mkdir(DataDir, os.ModePerm)

			if err != nil {
				panic(err)
			}
		}

		DBPath = path.Join(DataDir, "database.db")
		SecretPath = path.Join(DataDir, "secret.gob")
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
