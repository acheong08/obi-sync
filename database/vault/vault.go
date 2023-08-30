package vault

import (
	"errors"
	"log"
	"path"
	"time"

	"gorm.io/gorm"

	"github.com/acheong08/obsidian-sync/config"
	"github.com/acheong08/obsidian-sync/cryptography"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open(path.Join(config.DataDir, "vault.db")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&User{}, &Vault{}, &Share{})
	if err != nil {
		log.Fatal(err)
	}
}
func ShareVaultInvite(email, name, vaultID string) error {
	shareID := uuid.New().String()
	// _, err := db.DBConnection.Exec("INSERT INTO shares (id, email, name, vault_id) VALUES (?, ?, ?, ?)", shareID, email, name, vaultID)
	err := db.Create(&Share{
		UID:     shareID,
		Email:   email,
		Name:    name,
		VaultID: vaultID,
	}).Error
	return err
}

func ShareVaultRevoke(shareID, vaultID, userEmail string) error {
	var err error
	if shareID != "" {
		// _, err = db.DBConnection.Exec("DELETE FROM shares WHERE id = ? AND vault_id = ?", shareID, vaultID)
		err = db.Where("id = ? AND vault_id = ?", shareID, vaultID).Delete(&Share{}).Error
	} else {
		// _, err = db.DBConnection.Exec("DELETE FROM shares WHERE vault_id = ? AND email = ?", vaultID, userEmail)
		err = db.Where("vault_id = ? AND email = ?", vaultID, userEmail).Delete(&Share{}).Error
	}
	return err
}

func GetVaultShares(vaultID string) ([]*Share, error) {
	shares := []*Share{}
	err := db.Select("id, email, name, accepted").Where("vault_id = ?", vaultID).Find(&shares).Error
	return shares, err
}

func GetSharedVaults(userEmail string) ([]*Vault, error) {
	vaults := []*Vault{}
	vaultIDs := []string{}
	err := db.Model(&Share{}).Where("email = ?", userEmail).Select("vault_id").Scan(&vaultIDs).Error
	if err != nil {
		return nil, err
	}
	for _, vaultID := range vaultIDs {
		vault, err := GetVault(vaultID, "")
		if err != nil {
			return nil, err
		}
		vaults = append(vaults, vault)
	}
	return vaults, nil
}

func HasAccessToVault(vaultID, userEmail string) bool {
	if IsVaultOwner(vaultID, userEmail) {
		return true
	}

	// Check shares table
	var count int64
	// err := db.DBConnection.QueryRow("SELECT COUNT(*) FROM shares WHERE vault_id = ? AND email = ?", vaultID, userEmail).Scan(&count)
	err := db.Model(&Share{}).Where("vault_id = ? AND email = ?", vaultID, userEmail).Count(&count).Error
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return count > 0
}

func IsVaultOwner(vaultID, userEmail string) bool {
	var email string
	// Check vaults table
	// db.DBConnection.QueryRow("SELECT user_email FROM vaults WHERE id = ?", vaultID).Scan(&email)
	err := db.Model(&Vault{}).Where("id = ?", vaultID).Select("user_email").Scan(&email).Error
	return email == userEmail && err == nil
}

func NewUser(email, password, name string) error {
	// Create bcrypt hash of password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// _, err = db.DBConnection.Exec("INSERT INTO users (name, email, password, license) VALUES (?, ?, ?, ?)", name, email, hash, "")
	err = db.Create(&User{
		Name:     name,
		Email:    email,
		Password: string(hash),
		License:  "",
	}).Error
	return err
}
func UserInfo(email string) (*User, error) {
	var userInfo User
	// err := db.DBConnection.QueryRow("SELECT name, license FROM users WHERE email = ?", email).Scan(&name, &license)
	err := db.Model(&User{}).Where("email = ?", email).Select("name, license").Scan(&userInfo).Error
	return &userInfo, err
}
func Login(email, password string) (*User, error) {
	// Get user from database
	var user User
	// err := db.DBConnection.QueryRow("SELECT license, name, password FROM users WHERE email = ?", email).Scan(&user.License, &user.Name, &hash)
	err := db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	// Compare password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func NewVault(name, userEmail, password, salt, keyhash string) (*Vault, error) {
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
	newVault := &Vault{
		ID:        uuid.New().String(),
		UserEmail: userEmail,
		Created:   time.Now().UnixMilli(),
		Host:      config.Host,
		Name:      name,
		Password:  password,
		Salt:      salt,
		KeyHash:   keyhash,
	}
	err := db.Create(newVault).Error
	return newVault, err
}

func DeleteVault(id, email string) error {
	// _, err := db.DBConnection.Exec("DELETE FROM vaults WHERE id = ? AND user_email = ?", id, email)
	err := db.Where("id = ? AND user_email = ?", id, email).Delete(&Vault{}).Error
	return err
}

func GetVault(id, keyHash string) (*Vault, error) {
	vault := &Vault{}
	err := db.First(vault, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	if keyHash != "" {
		if vault.KeyHash != keyHash {
			return nil, errors.New("invalid keyhash")
		}
	}
	return vault, nil
}

func SetVaultVersion(id string, ver int) error {
	// _, err := db.DBConnection.Exec("UPDATE vaults SET version = ? WHERE id = ?", ver, id)
	err := db.Model(&Vault{}).Where("id = ?", id).Update("version", ver).Error
	return err
}

// Size is not included in the response. It should be fetched separately.
func GetVaults(userEmail string) ([]*Vault, error) {
	vaults := []*Vault{}
	err := db.Select("id, created, host, name, password, salt").Where("user_email = ?", userEmail).Find(&vaults).Error
	return vaults, err
}
