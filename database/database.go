package database

import (
	"database/sql"
	"os"

	"github.com/acheong08/obsidian-sync/config"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type Database struct {
	DBConnection *sql.DB
}

func NewDatabase() *Database {
	// Check if the database file exists
	if _, err := os.Stat(config.DBPath); os.IsNotExist(err) {
		db, err := sql.Open("sqlite", config.DBPath)
		if err != nil {
			panic(err)
		}
		// Create users table
		db.Exec(`CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL, 
			password TEXT NOT NULL,
			license TEXT,
			)`)
		// Create vaults table
		db.Exec(`CREATE TABLE vaults (
			id TEXT PRIMARY KEY,
			created INTEGER NOT NULL,
			host TEXT NOT NULL,
			name TEXT NOT NULL,
			password TEXT,
			salt TEXT NOT NULL,
			keyhash TEXT NOT NULL,
		)`)
		// Create files metadata table
		_, err = db.Exec(`CREATE TABLE files (
			vault_id TEXT NOT NULL,
			path TEXT NOT NULL,
			extension TEXT NOT NULL,
			hash TEXT NOT NULL,
			ctime INTEGER NOT NULL,
			mtime INTEGER NOT NULL,
			folder INTEGER NOT NULL,
			deleted INTEGER NOT NULL,
			size INTEGER NOT NULL,
			PRIMARY KEY (vault_id, path)
		)`)
		if err != nil {
			panic(err)
		}
		return &Database{
			DBConnection: db,
		}

	} else {
		// Connect to the database
		db, err := sql.Open("sqlite", config.DBPath)
		if err != nil {
			panic(err)
		}
		return &Database{
			DBConnection: db,
		}
	}
}

func (db *Database) Close() {
	db.DBConnection.Close()
}

func (db *Database) NewUser(email, password, name string) error {
	// Create bcrypt hash of password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.DBConnection.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", name, email, hash)
	return err
}
func (db *Database) Login(email, password string) (bool, error) {
	// Get user from database
	var hash string
	err := db.DBConnection.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&hash)
	if err != nil {
		return false, err
	}
	// Compare password hash
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil, err
}
