package publish

type Site struct {
	ID      string `json:"id,omitempty"`
	Host    string `json:"host,omitempty"`
	Created int64  `json:"created,omitempty"`
	Owner   string `json:"owner,omitempty"`
	Slug    string `json:"slug,omitempty"`
	Options string `json:"options,omitempty"`
	Size    int64  `json:"size"`
}
type File struct {
	Path  string `json:"path,omitempty"`
	CTime int64  `json:"ctime,omitempty"`
	Hash  string `json:"hash,omitempty"`
	MTime int64  `json:"mtime,omitempty"`
	Size  int64  `json:"size,omitempty"`
	Data  []byte `json:"data,omitempty"`
	Site  string `json:"site,omitempty"`
}
