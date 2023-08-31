package publish

import (
	"log"
	"path"
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
	db, err = gorm.Open(sqlite.Open(path.Join(config.DataDir, "publish.db")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&Site{}, &File{})
	if err != nil {
		log.Fatal(err)
	}
}

func GetFile(siteID, path string) (string, error) {
	var data string
	// err := db.QueryRow("SELECT data FROM files WHERE slug = ? AND path = ?", siteID, path).Scan(&data)
	err := db.Model(&File{}).Select("data").Where("slug = ? AND path = ?", siteID, path).First(&data).Error
	return data, err
}
func NewFile(file *File) error {
	file.CTime = time.Now().UnixMilli()
	file.MTime = time.Now().UnixMilli()
	// _, err := db.Exec("INSERT OR REPLACE INTO files (path, ctime, hash, mtime, size, data, site) VALUES (?, ?, ?, ?, ?, ?, ?)", file.Path, file.CTime, file.Hash, file.MTime, file.Size, file.Data, file.Slug)

	// Create with ON CONFLICT REPLACE
	err := db.Model(&File{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "path"}, {Name: "slug"}},
		UpdateAll: true,
	}).Create(file).Error
	return err
}

func RemoveFile(siteID, path string) error {
	// _, err := db.Exec("DELETE FROM files WHERE slug = ? AND path = ?", siteID, path)
	err := db.Model(&File{}).Where("slug = ? AND path = ?", siteID, path).Delete(&File{}).Error
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
	err := db.Create(&site).Error
	return &site, err
}

func DeleteSite(siteID string) error {
	err := db.Model(&Site{}).Where("id = ?", siteID).Delete(&Site{}).Error
	return err
}

type slugResponse struct {
	ID   string `json:"id"`
	Host string `json:"host"`
	Slug string `json:"slug"`
}

func GetSlug(slug string) (slugResponse, error) {
	// err := db.QueryRow("SELECT id, host FROM sites WHERE slug = ?", slug).Scan(&id, &host)
	var site slugResponse
	err := db.Model(&Site{}).Select("id, host, slug").Where("slug = ?", slug).First(&site).Error
	return site, err
}

func SetSlug(slug, id string) error {
	// _, err := db.Exec("UPDATE sites SET slug = ? WHERE id = ?", slug, id)
	err := db.Model(&Site{}).Where("id = ?", id).Update("slug", slug).Error
	return err
}

func GetSites(userEmail string) ([]*Site, error) {
	var sites []*Site = make([]*Site, 0)
	err := db.Model(&Site{}).Select("id, host, created, size").Where("owner = ?", userEmail).Find(&sites).Error
	return sites, err

}
func GetSiteOwner(siteID string) (string, error) {
	var email string
	// err := db.QueryRow("SELECT owner FROM sites WHERE id = ?", siteID).Scan(&email)
	err := db.Model(&Site{}).Select("owner").Where("id = ?", siteID).First(&email).Error
	return email, err
}

func GetSiteSlug(siteID string) (string, error) {
	var slug string
	// err := db.QueryRow("SELECT slug FROM sites WHERE id = ?", siteID).Scan(&slug)
	err := db.Model(&Site{}).Select("slug").Where("id = ?", siteID).First(&slug).Error
	return slug, err
}

func GetFiles(siteID string) ([]*File, error) {
	var files []*File = make([]*File, 0)
	err := db.Model(&File{}).Select("c_time, hash, m_time, path, size").Where("slug = ?", siteID).Find(&files).Error
	return files, err
}
