package database

import (
	"database/sql"
	"errors"
	"log"
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
		_, err = db.Exec(`CREATE TABLE users (
			name TEXT NOT NULL,
			email TEXT PRIMARY KEY NOT NULL,
			password TEXT NOT NULL,
			license TEXT NOT NULL
			)`)
		if err != nil {
			panic(err)
		}
		// Create vaults table
		_, err = db.Exec(`CREATE TABLE vaults (
			id TEXT PRIMARY KEY,
			user_email TEXT NOT NULL,
			created INTEGER NOT NULL,
			host TEXT NOT NULL,
			name TEXT NOT NULL,
			password TEXT NOT NULL,
			salt TEXT NOT NULL,
			version INTEGER NOT NULL DEFAULT 0,
			keyhash TEXT NOT NULL
		)`)
		if err != nil {
			panic(err)
		}
		// Create vault_shares table
		_, err = db.Exec(`CREATE TABLE shares (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			name TEXT NOT NULL,
			vault_id TEXT NOT NULL,
			accepted INTEGER NOT NULL DEFAULT 1
		)
		`)
		if err != nil {
			panic(err)
		}
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

func (db *Database) ShareVaultInvite(email, name, vaultID string) error {
	shareID := uuid.New().String()
	_, err := db.DBConnection.Exec("INSERT INTO shares (id, email, name, vault_id) VALUES (?, ?, ?, ?)", shareID, email, name, vaultID)
	return err
}

func (db *Database) ShareVaultRevoke(shareID, vaultID, userEmail string) error {
	var err error
	if shareID != "" {
		_, err = db.DBConnection.Exec("DELETE FROM shares WHERE id = ? AND vault_id = ?", shareID, vaultID)
	} else {
		_, err = db.DBConnection.Exec("DELETE FROM shares WHERE vault_id = ? AND email = ?", vaultID, userEmail)
	}
	return err
}

func (db *Database) GetVaultShares(vaultID string) ([]*vault.Share, error) {
	rows, err := db.DBConnection.Query("SELECT id, email, name, accepted FROM shares WHERE vault_id = ?", vaultID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	shares := []*vault.Share{}
	for rows.Next() {
		share := &vault.Share{}
		err = rows.Scan(&share.UID, &share.Email, &share.Name, &share.Accepted)
		if err != nil {
			return nil, err
		}
		shares = append(shares, share)
	}
	return shares, nil
}

func (db *Database) GetSharedVaults(userEmail string) ([]*vault.Vault, error) {
	rows, err := db.DBConnection.Query("SELECT vault_id FROM shares WHERE email = ?", userEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	vaults := []*vault.Vault{}
	for rows.Next() {
		var vaultID string
		err = rows.Scan(&vaultID)
		if err != nil {
			return nil, err
		}
		vault, err := db.GetVault(vaultID, "")
		if err != nil {
			return nil, err
		}
		vaults = append(vaults, vault)
	}
	return vaults, nil
}

func (db *Database) HasAccessToVault(vaultID, userEmail string) bool {
	if db.IsVaultOwner(vaultID, userEmail) {
		return true
	}

	// Check shares table
	var count int
	err := db.DBConnection.QueryRow("SELECT COUNT(*) FROM shares WHERE vault_id = ? AND email = ?", vaultID, userEmail).Scan(&count)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func (db *Database) IsVaultOwner(vaultID, userEmail string) bool {
	var email string
	// Check vaults table
	db.DBConnection.QueryRow("SELECT user_email FROM vaults WHERE id = ?", vaultID).Scan(&email)
	return email == userEmail
}

func (db *Database) NewUser(email, password, name string) error {
	// Create bcrypt hash of password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.DBConnection.Exec("INSERT INTO users (name, email, password, license) VALUES (?, ?, ?, ?)", name, email, hash, "")
	return err
}
func (db *Database) UserInfo(email string) (*user.User, error) {
	var name string
	var license string
	err := db.DBConnection.QueryRow("SELECT name, license FROM users WHERE email = ?", email).Scan(&name, &license)
	if err != nil {
		return nil, err
	}
	return &user.User{
		Name:    name,
		License: license,
	}, nil
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

func (db *Database) NewVault(name, userEmail, password, salt, keyhash string) (*vault.Vault, error) {
	if keyhash == "" && password == "" {
		return nil, errors.New("password and keyhash cannot both be empty")
	}
	if keyhash == "" {
		var err error
		keyhash, err = cryptography.MakeKeyHash(password, salt)
		if err != nil {
			return nil, err
		}
	}
	id := uuid.New().String()
	created := time.Now().UnixMilli()
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
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, id, userEmail, created, host, name, password, salt, keyhash)
	return &vault.Vault{
		ID:       id,
		Created:  created,
		Host:     host,
		Name:     name,
		Password: password,
		Salt:     salt,
		Size:     0,
	}, err
}

func (db *Database) DeleteVault(id, email string) error {
	_, err := db.DBConnection.Exec("DELETE FROM vaults WHERE id = ? AND user_email = ?", id, email)
	return err
}

func (db *Database) GetVault(id, keyHash string) (*vault.Vault, error) {
	vault := &vault.Vault{}
	var dbKeyHash string
	err := db.DBConnection.QueryRow("SELECT * FROM vaults WHERE id = ?", id).Scan(
		&vault.ID, &vault.UserEmail, &vault.Created, &vault.Host, &vault.Name, &vault.Password, &vault.Salt,
		&vault.Version, &dbKeyHash)
	if err != nil {
		return nil, err
	}
	if keyHash != "" {
		if dbKeyHash != keyHash {
			return nil, errors.New("invalid keyhash")
		}
	}
	return vault, nil
}

func (db *Database) SetVaultVersion(id string, ver int) error {
	_, err := db.DBConnection.Exec("UPDATE vaults SET version = ? WHERE id = ?", ver, id)
	return err
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
	if err != nil {
		return nil, err
	}
	return vaults, nil
}
