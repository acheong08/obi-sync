package vault

type Vault struct {
	ID       string `json:"id"`
	Created  int64  `json:"created"`
	Host     string `json:"host"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Size     int64  `json:"size"`
	// KeyHash  string `json:"keyhash"`
}
