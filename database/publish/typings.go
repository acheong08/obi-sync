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
	Path    string `json:"path,omitempty" gorm:"uniqueIndex:idx_path_site"`
	CTime   int64  `json:"ctime,omitempty"`
	Hash    string `json:"hash,omitempty"`
	MTime   int64  `json:"mtime,omitempty"`
	Size    int64  `json:"size,omitempty"`
	Data    string `json:"data,omitempty"`
	Slug    string `json:"site,omitempty" gorm:"uniqueIndex:idx_path_site"`
	Deleted bool   `json:"deleted"`
}
