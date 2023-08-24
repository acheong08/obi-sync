package config

import (
	"crypto/rand"
	"encoding/gob"
	"os"
)

const DBPath = "database.db"
const Host = "localhost:3000/ws"

var Secret []byte

// Generate a random password, hash it, and store it in the Secret variable & a file
// Load secret.gob if it exists
func init() {
	if _, err := os.Stat("secret.gob"); err != nil {
		Secret = make([]byte, 64)
		rand.Read(Secret)
		f, err := os.Create("secret.gob")
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
		f, err := os.Open("secret.gob")
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
