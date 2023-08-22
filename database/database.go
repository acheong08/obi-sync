package database

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/acheong08/obsidian-sync/user"

	"github.com/acheong08/obsidian-sync/config"
	"github.com/acheong08/obsidian-sync/cryptography"
	"github.com/acheong08/obsidian-sync/vault"
	"github.com/google/uuid"
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
			name TEXT NOT NULL,
			email TEXT PRIMARY KEY NOT NULL,
			password TEXT NOT NULL,
			license TEXT,
			)`)
		// Create vaults table
		db.Exec(`CREATE TABLE vaults (
			id TEXT PRIMARY KEY,
			user_email TEXT NOT NULL,
			created INTEGER NOT NULL,
			host TEXT NOT NULL,
			name TEXT NOT NULL,
			password TEXT NOT NULL,
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
func (db *Database) Login(email, password string) (*user.User, error) {
	// Get user from database
	var hash string
	var user user.User
	err := db.DBConnection.QueryRow("SELECT license, name, password FROM users WHERE email = ?", email).Scan(&user.License, &user.Name, &hash)
	if err != nil {
		return nil, err
	}
	// Compare password hash
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return nil, err
	}
	user.Email = email
	return &user, nil
}

func (db *Database) NewVault(name, userEmail, password, salt, keyhash string) error {
	if keyhash == "" && password == "" {
		return errors.New("password and keyhash cannot both be empty")
	}
	if keyhash == "" {
		var err error
		keyhash, err = cryptography.MakeKeyHash(password, salt)
		if err != nil {
			return err
		}
	}
	id := uuid.New().String()
	created := time.Now().Unix()
	host := config.Host
	_, err := db.DBConnection.Exec(`INSERT INTO vaults (
			id,
			user_email,
			created,
			host,
			name, 
			password, 
			salt, 
			keyhash
		) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`, id, userEmail, created, host, name, password, salt, keyhash)
	return err
}

func (db *Database) GetVault(id, keyHash string) (*vault.Vault, error) {
	vault := &vault.Vault{}
	var dbKeyHash string
	err := db.DBConnection.QueryRow("SELECT * FROM vaults WHERE id = ?", id).Scan(&vault.ID, &vault.Created, &vault.Host, &vault.Name, &vault.Password, &vault.Salt, &dbKeyHash)
	if err != nil {
		return nil, err
	}
	if dbKeyHash != keyHash {
		return nil, errors.New("invalid keyhash")
	}
	return vault, nil
}

// Size is not included in the response. It should be fetched separately.
func (db *Database) GetVaults(userEmail string) ([]*vault.Vault, error) {
	rows, err := db.DBConnection.Query("SELECT id, created, host, name, password, salt FROM vaults WHERE user_email = ?", userEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	vaults := []*vault.Vault{}
	for rows.Next() {
		vault := &vault.Vault{}
		err = rows.Scan(&vault.ID, &vault.Created, &vault.Host, &vault.Name, &vault.Password, &vault.Salt)
		if err != nil {
			return nil, err
		}
		vaults = append(vaults, vault)
	}
	return vaults, nil
}
