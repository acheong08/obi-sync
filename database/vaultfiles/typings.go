package vaultfiles

type File struct {
	UID        int    `json:"uid" gorm:"primary_key;autoIncrement"`
	VaultID    string `json:"vault_id"`
	Hash       string `json:"hash"`
	Path       string `json:"path"`
	Extension  string `json:"extension"`
	Size       int64  `json:"size"`
	Created    int64  `json:"created"`
	Modified   int64  `json:"modified"`
	Folder     bool   `json:"folder"`
	Deleted    bool   `json:"deleted"`
	Data       []byte `json:"-"`
	Newest     bool   `json:"newest" gorm:"default:true"`
	IsSnapshot bool   `json:"is_snapshot" gorm:"default:false"`
}

// 	CREATE TABLE IF NOT EXISTS files (
// 		uid INTEGER PRIMARY KEY AUTOINCREMENT,
// 		vault_id TEXT,
// 		hash TEXT,
// 		path TEXT,
// 		extension TEXT,
// 		size INTEGER,
// 		created INTEGER,
// 		modified INTEGER,
// 		folder INTEGER,
// 		deleted INTEGER,
// 		data BLOB,
// 		newest INTEGER NOT NULL DEFAULT 1,
// 		is_snapshot INTEGER NOT NULL DEFAULT 0

type FileResponse struct {
	File
	OP string `json:"op"`
}
