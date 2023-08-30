package vaultfiles

type File struct {
	UID        int    `json:"uid" gorm:"primary_key"`
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
	Newest     bool   `json:"newest"`
	IsSnapshot bool   `json:"is_snapshot"`
}

type FileResponse struct {
	File
	OP string `json:"op"`
}
