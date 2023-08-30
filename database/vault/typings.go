package vault

type Vault struct {
	ID        string `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	UserEmail string `json:"user_email,omitempty"`
	Created   int64  `json:"created,omitempty"`
	Host      string `json:"host,omitempty"`
	Name      string `json:"name,omitempty"`
	Password  string `json:"password,omitempty"`
	Salt      string `json:"salt,omitempty"`
	Size      int64  `json:"size,omitempty"`
	// Not part of JSON
	Version int    `json:"-,omitempty"`
	KeyHash string `json:"keyhash,omitempty"`
}

type Share struct {
	UID      string `json:"uid,omitempty"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	VaultID  string `json:"vault_id,omitempty"`
	Accepted bool   `json:"accepted,omitempty"`
}
