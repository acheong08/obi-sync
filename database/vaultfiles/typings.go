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
	Folder     int    `json:"folder"`
	Deleted    int    `json:"deleted"`
	Data       []byte `json:"-"`
	Newest     int    `json:"newest"`
	IsSnapshot int    `json:"is_snapshot"`
}

type FileResponse struct {
	File
	OP string `json:"op"`
}
