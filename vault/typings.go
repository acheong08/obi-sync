package vault

type Vault struct {
	ID        string `json:"id"`
	UserEmail string `json:"user_email,omitempty"`
	Created   int64  `json:"created"`
	Host      string `json:"host"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Salt      string `json:"salt"`
	Size      int64  `json:"size"`
	// Not part of JSON
	Version int `json:"-"`
	// KeyHash  string `json:"keyhash"`
}

type File struct {
	Op          string `json:"op,omitempty"`
	Path        string `json:"path"`
	Hash        string `json:"hash,omitempty"`
	Extension   string `json:"extension,omitempty"`
	Size        int64  `json:"size"`
	Created     int64  `json:"ctime,omitempty"`
	Modified    int64  `json:"mtime,omitempty"`
	Timestamp   int64  `json:"ts,omitempty"`
	RelatedPath string `json:"related_path,omitempty"`
	Folder      bool   `json:"folder"`
	Deleted     bool   `json:"deleted"`
	UID         int    `json:"uid"`
	Data        []byte `json:"-"`
}
