package publish

import (
	"database/sql"
	"log"
	"path"
	"time"

	"github.com/acheong08/obsidian-sync/config"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func init() {
	var err error
	var dbPath = path.Join(config.DataDir, "publish.db")
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sites (
			id TEXT PRIMARY KEY,
			host TEXT NOT NULL,
			created INTEGER NOT NULL,
			owner TEXT NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			options TEXT,
			size INTEGER NOT NULL DEFAULT 0
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			path TEXT PRIMARY KEY,
			ctime INTEGER NOT NULL,
			hash TEXT NOT NULL,
			mtime INTEGER NOT NULL,
			size INTEGER NOT NULL,
			data BLOB NOT NULL,
			site TEXT NOT NULL,
			deleted INTEGER NOT NULL DEFAULT 0,
			UNIQUE (path, site)
		);`)
	if err != nil {
		log.Fatal(err)
	}
}

func GetFile(siteID, path string) ([]byte, error) {
	var data []byte
	err := db.QueryRow("SELECT data FROM files WHERE site = ? AND path = ?", siteID, path).Scan(&data)
	return data, err
}
func NewFile(file *File) error {
	file.CTime = time.Now().UnixMilli()
	file.MTime = time.Now().UnixMilli()
	_, err := db.Exec("INSERT OR REPLACE INTO files (path, ctime, hash, mtime, size, data, site) VALUES (?, ?, ?, ?, ?, ?, ?)", file.Path, file.CTime, file.Hash, file.MTime, file.Size, file.Data, file.Site)
	return err
}

func RemoveFile(siteID, path string) error {
	_, err := db.Exec("DELETE FROM files WHERE site = ? AND path = ?", siteID, path)
	return err
}

func CreateSite(owner string) (*Site, error) {
	var site Site = Site{
		ID:      uuid.New().String(),
		Host:    config.Host,
		Created: time.Now().UnixMilli(),
		Owner:   owner,
		Slug:    uuid.New().String(),
	}
	_, err := db.Exec("INSERT INTO sites (id, host, created, owner, slug) VALUES (?, ?, ?, ?, ?)", site.ID, site.Host, site.Created, site.Owner, site.Slug)
	return &site, err
}

type slugResponse struct {
	ID   string `json:"id"`
	Host string `json:"host"`
	Slug string `json:"slug"`
}

func GetSlug(slug string) (slugResponse, error) {
	var (
		id, host string
	)
	err := db.QueryRow("SELECT id, host FROM sites WHERE slug = ?", slug).Scan(&id, &host)
	return slugResponse{
		ID:   id,
		Host: host,
		Slug: slug,
	}, err
}

func SetSlug(slug, id string) error {
	_, err := db.Exec("UPDATE sites SET slug = ? WHERE id = ?", slug, id)
	return err
}

func GetSites(userEmail string) ([]*Site, error) {
	var sites []*Site = make([]*Site, 0)
	rows, err := db.Query("SELECT id, host, created, size FROM sites WHERE owner = ?", userEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var site Site
		err := rows.Scan(&site.ID, &site.Host, &site.Created, &site.Size)
		if err != nil {
			return nil, err
		}
		sites = append(sites, &site)
	}
	return sites, nil

}
func GetSiteOwner(siteID string) (string, error) {
	var email string
	err := db.QueryRow("SELECT owner FROM sites WHERE id = ?", siteID).Scan(&email)
	return email, err
}

func GetSiteSlug(siteID string) (string, error) {
	var slug string
	err := db.QueryRow("SELECT slug FROM sites WHERE id = ?", siteID).Scan(&slug)
	return slug, err
}

func GetFiles(siteID string) ([]*File, error) {
	var files []*File = make([]*File, 0)
	rows, err := db.Query("SELECT ctime, hash, mtime, path, size FROM files WHERE site = ?", siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var file File
		err := rows.Scan(&file.CTime, &file.Hash, &file.MTime, &file.Path, &file.Size)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}
	return files, nil
}
