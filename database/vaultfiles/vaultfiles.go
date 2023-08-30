package vaultfiles

import (
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/glebarez/sqlite"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("vault.db"), &gorm.Config{})
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
	// _, err := db.Exec(`
	// 	UPDATE files SET is_snapshot = 1 WHERE newest = 1 AND vault_id = $1;
	// 	DELETE FROM files WHERE is_snapshot = 0 AND vault_id = $1;
	// 	DELETE FROM files WHERE size != 0 AND data IS NULL AND vault_id = $1;
	// `, vaultID)

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
	// err := db.QueryRow("SELECT path, hash, extension, size, created, modified, folder, deleted FROM files WHERE uid = ?", uid).Scan(&file.Path, &file.Hash, &file.Extension, &file.Size, &file.Created, &file.Modified, &file.Folder, &file.Deleted)
	err := db.Select("path, hash, extension, size, created, modified, folder, deleted").Where("uid = ?", uid).First(&file).Error
	if err != nil {
		return nil, err
	}
	file.UID = uid
	// Update file to be not deleted and newest and all other files with the
	// same path to not be newest
	// _, err = db.Exec(`UPDATE files SET deleted = 0, newest = 1 WHERE uid = $1;
	// 	UPDATE files SET newest = 0 WHERE path = $2 AND deleted = 0`, uid, file.Path)
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
	// rows, err := db.Query("SELECT uid, path, hash, extension, size, created, modified, folder, deleted FROM files WHERE vault_id = ? AND deleted = 0 AND newest = 1", vaultID)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()
	//
	// for rows.Next() {
	// 	var file File
	// 	err = rows.Scan(&file.UID, &file.Path, &file.Hash, &file.Extension, &file.Size, &file.Created, &file.Modified, &file.Folder, &file.Deleted)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	files = append(files, file)
	// }
	err := db.Model(&File{}).Select("uid, path, hash, extension, size, created, modified, folder, deleted").Where("vault_id = ? AND deleted = 0 AND newest = 1", vaultID).Find(&files).Error
	return files, err
}

func GetFile(uid int) (*File, error) {
	var file File
	// Get hash and size
	// err := db.QueryRow("SELECT hash, size, data FROM files WHERE uid = ?", uid).Scan(&file.Hash, &file.Size, &file.Data)
	err := db.Model(&File{}).Select("hash, size, data").Where("uid = ?", uid).First(&file).Error
	return &file, err
}

func GetFileHistory(path string) ([]*File, error) {
	var files []*File
	// Order by modified time (newest first in array)
	// rows, err := db.Query("SELECT uid, path, size, modified, folder, deleted FROM files WHERE path = ? ORDER BY modified DESC", path)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var file File
	// 	err = rows.Scan(&file.UID, &file.Path, &file.Size, &file.Timestamp, &file.Folder, &file.Deleted)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	files = append(files, file)
	// }
	err := db.Model(&File{}).Select("uid, path, size, modified, folder, deleted").Where("path = ?", path).Order("modified DESC").Find(&files).Error
	return files, err
}

func GetDeletedFiles() (any, error) {
	type f struct {
		UID       int    `json:"uid"`
		Timestamp int64  `json:"ts"`
		Size      int64  `json:"size"`
		Path      string `json:"path"`
		Folder    bool   `json:"folder"`
		Deleted   bool   `json:"deleted"`
	}
	var files []*f = make([]*f, 0)
	// Get all files that are deleted (deleted,folder,path,size,modified,uid)
	// rows, err := db.Query("SELECT uid, modified, size, path, folder, deleted FROM files WHERE deleted = ?", true)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var file f
	// 	err = rows.Scan(&file.UID, &file.Timestamp, &file.Size, &file.Path, &file.Folder, &file.Deleted)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	files = append(files, file)

	// }
	// if files == nil {
	// 	return make([]File, 0), nil
	// }
	err := db.Model(&File{}).Select("uid, modified, size, path, folder, deleted").Where("deleted = ?", true).Find(&files).Error
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
	// Set previous files with the same path to not be newest
	// _, err := db.Exec("UPDATE files SET newest = 0 WHERE path = ? AND newest = 1", file.Path)
	err := db.Model(&File{}).Where("path = ? AND newest = 1", file.Path).Update("newest", 0).Error
	if err != nil {
		return 0, err
	}
	// result, err := db.Exec(`INSERT INTO files (
	// 	vault_id, path, hash, extension, size, created, modified, folder, deleted)
	// 	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
	// 	vaultID, file.Path, file.Hash, file.Extension, file.Size, file.Created,
	// 	file.Modified, file.Folder, file.Deleted)
	result := db.Create(file)
	if result.Error != nil {
		return 0, result.Error
	}

	return file.UID, err
}

func GetFileData(uid int) (*[]byte, error) {
	var file []byte
	// err := db.QueryRow("SELECT data FROM files WHERE uid = ?", uid).Scan(&file)
	err := db.Select("data").Where("uid = ?", uid).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func InsertData(uid int, data *[]byte) error {
	// _, err := db.Exec("UPDATE files SET data = ? WHERE uid = ?", data, uid)
	err := db.Model(&File{}).Where("uid = ?", uid).Update("data", data).Error
	return err
}

func DeleteVaultFile(path string) error {
	// Update all files with the same path to be deleted
	// _, err := db.Exec("UPDATE files SET deleted = 1 WHERE path = ?", path)
	err := db.Model(&File{}).Where("path = ?", path).Update("deleted", 1).Error
	return err
}
