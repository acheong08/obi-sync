package publish

import (
	"log"
	"time"

	"github.com/acheong08/obsidian-sync/config"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("publish.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&Site{}, &File{})
	if err != nil {
		log.Fatal(err)
	}
}

func GetFile(siteID, path string) ([]byte, error) {
	var data []byte
	// err := db.QueryRow("SELECT data FROM files WHERE site = ? AND path = ?", siteID, path).Scan(&data)
	err := db.Select("data").Where("site = ? AND path = ?", siteID, path).First(&data).Error
	return data, err
}
func NewFile(file *File) error {
	file.CTime = time.Now().UnixMilli()
	file.MTime = time.Now().UnixMilli()
	// _, err := db.Exec("INSERT OR REPLACE INTO files (path, ctime, hash, mtime, size, data, site) VALUES (?, ?, ?, ?, ?, ?, ?)", file.Path, file.CTime, file.Hash, file.MTime, file.Size, file.Data, file.Slug)

	// Create with ON CONFLICT REPLACE
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "path"}, {Name: "site"}},
		UpdateAll: true,
	}).Create(file).Error
	return err
}

func RemoveFile(siteID, path string) error {
	// _, err := db.Exec("DELETE FROM files WHERE site = ? AND path = ?", siteID, path)
	err := db.Where("site = ? AND path = ?", siteID, path).Delete(&File{}).Error
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
	// _, err := db.Exec("INSERT INTO sites (id, host, created, owner, slug) VALUES (?, ?, ?, ?, ?)", site.ID, site.Host, site.Created, site.Owner, site.Slug)
	err := db.Create(&site).Error
	return &site, err
}

type slugResponse struct {
	ID   string `json:"id"`
	Host string `json:"host"`
	Slug string `json:"slug"`
}

func GetSlug(slug string) (slugResponse, error) {
	// err := db.QueryRow("SELECT id, host FROM sites WHERE slug = ?", slug).Scan(&id, &host)
	var site slugResponse
	err := db.Select("id, host, slug").Where("slug = ?", slug).First(&site).Error
	return site, err
}

func SetSlug(slug, id string) error {
	// _, err := db.Exec("UPDATE sites SET slug = ? WHERE id = ?", slug, id)
	err := db.Model(&Site{}).Where("id = ?", id).Update("slug", slug).Error
	return err
}

func GetSites(userEmail string) ([]*Site, error) {
	var sites []*Site = make([]*Site, 0)
	// rows, err := db.Query("SELECT id, host, created, size FROM sites WHERE owner = ?", userEmail)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var site Site
	// 	err := rows.Scan(&site.ID, &site.Host, &site.Created, &site.Size)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	sites = append(sites, &site)
	// }
	err := db.Select("id, host, created, size").Where("owner = ?", userEmail).Find(&sites).Error
	return sites, err

}
func GetSiteOwner(siteID string) (string, error) {
	var email string
	// err := db.QueryRow("SELECT owner FROM sites WHERE id = ?", siteID).Scan(&email)
	err := db.Select("owner").Where("id = ?", siteID).First(&email).Error
	return email, err
}

func GetSiteSlug(siteID string) (string, error) {
	var slug string
	// err := db.QueryRow("SELECT slug FROM sites WHERE id = ?", siteID).Scan(&slug)
	err := db.Select("slug").Where("id = ?", siteID).First(&slug).Error
	return slug, err
}

func GetFiles(siteID string) ([]*File, error) {
	var files []*File = make([]*File, 0)
	// rows, err := db.Query("SELECT ctime, hash, mtime, path, size FROM files WHERE site = ?", siteID)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var file File
	// 	err := rows.Scan(&file.CTime, &file.Hash, &file.MTime, &file.Path, &file.Size)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	files = append(files, &file)
	// }
	err := db.Select("ctime, hash, mtime, path, size").Where("site = ?", siteID).Find(&files).Error
	return files, err
}
