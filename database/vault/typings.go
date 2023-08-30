package vault

type Vault struct {
	ID        string `json:"id,omitempty" gorm:"primaryKey"`
	UserEmail string `json:"user_email,omitempty"`
	Created   int64  `json:"created,omitempty"`
	Host      string `json:"host,omitempty"`
	Name      string `json:"name,omitempty"`
	Password  string `json:"password,omitempty"`
	Salt      string `json:"salt,omitempty"`
	Size      int64  `json:"size,omitempty"`
	// Not part of JSON
	Version int    `json:"-,omitempty" gorm:"default:0"`
	KeyHash string `json:"keyhash,omitempty"`
}

// CREATE TABLE vaults (
// 		id TEXT PRIMARY KEY,
// 		user_email TEXT NOT NULL,
// 		created INTEGER NOT NULL,
// 		host TEXT NOT NULL,
// 		name TEXT NOT NULL,
// 		password TEXT NOT NULL,
// 		salt TEXT NOT NULL,
// 		version INTEGER NOT NULL DEFAULT 0,
// 		keyhash TEXT NOT NULL
// 	)

type Share struct {
	UID      string `json:"uid,omitempty" gorm:"primaryKey"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	VaultID  string `json:"vault_id,omitempty"`
	Accepted bool   `json:"accepted,omitempty" gorm:"default:true"`
}

// CREATE TABLE shares (
// 		id TEXT PRIMARY KEY,
// 		email TEXT NOT NULL,
// 		name TEXT NOT NULL,
// 		vault_id TEXT NOT NULL,
// 		accepted INTEGER NOT NULL DEFAULT 1
// 	)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"primaryKey"`
	Password string `json:"password"`
	License  string `json:"license"`
}

// CREATE TABLE users (
// name TEXT NOT NULL,
// email TEXT PRIMARY KEY NOT NULL,
// password TEXT NOT NULL,
// license TEXT NOT NULL
// )
