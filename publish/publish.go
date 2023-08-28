package publish

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite", "publish.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sites (
			id TEXT PRIMARY KEY,
			host TEXT NOT NULL,
			created INTEGER NOT NULL,
			owner TEXT NOT NULL,
			slug TEXT NOT NULL,
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
			data TEXT NOT NULL,
			site TEXT NOT NULL,
			deleted INTEGER NOT NULL DEFAULT 0
		);`)
	if err != nil {
		log.Fatal(err)
	}

}

func GetSites(userEmail string) ([]*Site, error) {
	var sites []*Site
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
	var files []*File
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
