package vault

type Vault struct {
	ID       string `json:"id"`
	Created  int64  `json:"created"`
	Host     string `json:"host"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Size     int64  `json:"size"`
	// Not part of JSON
	Version int `json:"-"`
	// KeyHash  string `json:"keyhash"`
}

type FileMetadata struct {
	Path     string `json:"path"`
	Hash     string `json:"hash"`
	Size     int64  `json:"size"`
	Created  int64  `json:"ctime"`
	Modified int64  `json:"mtime"`
	Folder   bool   `json:"folder"`
	Deleted  bool   `json:"deleted"`
	UID      string `json:"uid"`
}
