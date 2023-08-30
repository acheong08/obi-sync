package vaultfiles

import (
	"log"
	"path"
	"time"

	"gorm.io/gorm"

	"github.com/acheong08/obsidian-sync/config"
	"github.com/glebarez/sqlite"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open(path.Join(config.DataDir, "vaultfiles.db")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&File{})
	if err != nil {
		log.Fatal(err)
	}
}

// Sets the newest files to also be snapshots and
// deletes all files which are not snapshots
// deletes all files where size is not 0 but data is null
func Snapshot(vaultID string) error {
	// Set newest files to be snapshots
	err := db.Model(&File{}).Where("newest = 1 AND vault_id = ?", vaultID).Update("is_snapshot", 1).Error
	if err != nil {
		return err
	}
	// Delete all files which are not snapshots
	err = db.Where("is_snapshot = 0 AND vault_id = ?", vaultID).Delete(&File{}).Error
	if err != nil {
		return err
	}
	// Delete all files where size is not 0 but data is null
	err = db.Where("size != 0 AND data IS NULL AND vault_id = ?", vaultID).Delete(&File{}).Error
	return err
}

func RestoreFile(uid int) (*FileResponse, error) {
	// Get file path
	var file File
	err := db.Select("path, hash, extension, size, created, modified, folder, deleted").Where("uid = ?", uid).First(&file).Error
	if err != nil {
		return nil, err
	}
	file.UID = uid
	err = db.Model(&File{}).Where("uid = ?", uid).Updates(File{
		Deleted: false,
		Newest:  true,
	}).Error
	if err != nil {
		return nil, err
	}
	err = db.Model(&File{}).Where("path = ? AND deleted = 0", file.Path).Update("newest", 0).Error
	return &FileResponse{
		file,
		"push",
	}, err
}

func GetVaultSize(vaultID string) (int, error) {
	var size int
	// err := db.QueryRow("SELECT COALESCE(SUM(size), 0) FROM files WHERE vault_id = ?", vaultID).Scan(&size)
	err := db.Model(&File{}).Select("COALESCE(SUM(size), 0)").Where("vault_id = ?", vaultID).First(&size).Error
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return size, nil
}

func GetVaultFiles(vaultID string) ([]*File, error) {
	var files []*File
	err := db.Model(&File{}).Select("uid, path, hash, extension, size, created, modified, folder, deleted").Where("vault_id = ? AND deleted = 0 AND newest = 1", vaultID).Find(&files).Error
	return files, err
}

func GetFile(uid int) (*File, error) {
	var file File
	err := db.Model(&File{}).Select("hash, size, data").Where("uid = ?", uid).First(&file).Error
	return &file, err
}

func GetFileHistory(path string) ([]*File, error) {
	var files []*File
	err := db.Model(&File{}).Select("uid, path, size, modified, folder, deleted").Where("path = ?", path).Order("modified DESC").Find(&files).Error
	return files, err
}

func GetDeletedFiles() (any, error) {
	type f struct {
		UID      int    `json:"uid"`
		Modified int64  `json:"ts"`
		Size     int64  `json:"size"`
		Path     string `json:"path"`
		Folder   bool   `json:"folder"`
		Deleted  bool   `json:"deleted"`
	}
	var files []*f = make([]*f, 0)
	// err := db.Model(&File{}).Select("uid, modified, size, path, folder, deleted").Where("deleted = ?", true).Find(&files).Error
	// Find all files which are deleted and newest
	err := db.Model(&File{}).Select("uid, modified, size, path, folder, deleted").Where("deleted = ? AND newest = ?", true, true).Find(&files).Error
	return &files, err
}

func InsertMetadata(file *File) (int, error) {
	// If created & modified are 0, set them to current time
	if file.Created == 0 {
		file.Created = time.Now().UnixMilli()
	}
	if file.Modified == 0 {
		file.Modified = time.Now().UnixMilli()
	}
	err := db.Model(&File{}).Where("path = ? AND newest = 1", file.Path).Update("newest", 0).Error
	if err != nil {
		return 0, err
	}
	result := db.Create(file)
	if result.Error != nil {
		return 0, result.Error
	}

	return file.UID, err
}

func InsertData(uid int, data *[]byte) error {
	err := db.Model(&File{}).Where("uid = ?", uid).Update("data", data).Error
	return err
}

func DeleteVaultFile(path string) error {
	// Set deleted to true and is_snapshot to true
	err := db.Model(&File{}).Where("path = ?", path).Updates(File{
		Deleted:    true,
		IsSnapshot: true,
	}).Error
	return err
}
