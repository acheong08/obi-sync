package vault

import (
	"database/sql"
	"log"
	"path"
	"time"

	"github.com/acheong08/obsidian-sync/config"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func init() {
	var err error
	var dbPath = path.Join(config.DataDir, "vault.db")
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create tables if they don't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			uid INTEGER PRIMARY KEY AUTOINCREMENT,
			vault_id TEXT,
			hash TEXT,
			path TEXT,
			extension TEXT,
			size INTEGER,
			created INTEGER,
			modified INTEGER,
			folder INTEGER,
			deleted INTEGER,
			data BLOB,
			newest INTEGER NOT NULL DEFAULT 1,
			is_snapshot INTEGER NOT NULL DEFAULT 0
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

// Sets the newest files to also be snapshots and
// deletes all files which are not snapshots
// deletes all files where size is not 0 but data is null
func Snapshot(vaultID string) error {
	_, err := db.Exec(`
		UPDATE files SET is_snapshot = 1 WHERE newest = 1 AND vault_id = $1;
		DELETE FROM files WHERE is_snapshot = 0 AND vault_id = $1;
		DELETE FROM files WHERE size != 0 AND data IS NULL AND vault_id = $1;
	`, vaultID)
	return err
}

func RestoreFile(uid int) (*File, error) {
	// Get file path
	var file File
	err := db.QueryRow("SELECT path, hash, extension, size, created, modified, folder, deleted FROM files WHERE uid = ?", uid).Scan(&file.Path, &file.Hash, &file.Extension, &file.Size, &file.Created, &file.Modified, &file.Folder, &file.Deleted)
	if err != nil {
		return nil, err
	}
	file.UID = uid
	// Update file to be not deleted and newest and all other files with the
	// same path to not be newest
	_, err = db.Exec(`UPDATE files SET deleted = 0, newest = 1 WHERE uid = $1;
		UPDATE files SET newest = 0 WHERE path = $2 AND deleted = 0`, uid, file.Path)

	return &file, err
}

func GetVaultSize(vaultID string) (int64, error) {
	var size sql.NullInt64
	err := db.QueryRow("SELECT COALESCE(SUM(size), 0) FROM files WHERE vault_id = ?", vaultID).Scan(&size)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	if size.Valid {
		return size.Int64, nil
	}
	return 0, nil
}

func GetVaultFiles(vaultID string) (*[]File, error) {
	rows, err := db.Query("SELECT uid, path, hash, extension, size, created, modified, folder, deleted FROM files WHERE vault_id = ? AND deleted = 0 AND newest = 1", vaultID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var files []File
	for rows.Next() {
		var file File
		err = rows.Scan(&file.UID, &file.Path, &file.Hash, &file.Extension, &file.Size, &file.Created, &file.Modified, &file.Folder, &file.Deleted)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return &files, nil
}

func GetFile(uid int) (*File, error) {
	var file File
	// Get hash and size
	err := db.QueryRow("SELECT hash, size, data FROM files WHERE uid = ?", uid).Scan(&file.Hash, &file.Size, &file.Data)
	return &file, err
}

func GetFileHistory(path string) (*[]File, error) {
	// Order by modified time (newest first in array)
	rows, err := db.Query("SELECT uid, path, size, modified, folder, deleted FROM files WHERE path = ? ORDER BY modified DESC", path)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var files []File
	for rows.Next() {
		var file File
		err = rows.Scan(&file.UID, &file.Path, &file.Size, &file.Timestamp, &file.Folder, &file.Deleted)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return &files, nil
}

func GetDeletedFiles() (any, error) {
	// Get all files that are deleted (deleted,folder,path,size,modified,uid)
	rows, err := db.Query("SELECT uid, modified, size, path, folder, deleted FROM files WHERE deleted = ?", true)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	type f struct {
		UID       int    `json:"uid"`
		Timestamp int64  `json:"ts"`
		Size      int64  `json:"size"`
		Path      string `json:"path"`
		Folder    bool   `json:"folder"`
		Deleted   bool   `json:"deleted"`
	}
	var files []f
	for rows.Next() {
		var file f
		err = rows.Scan(&file.UID, &file.Timestamp, &file.Size, &file.Path, &file.Folder, &file.Deleted)
		if err != nil {
			return nil, err
		}
		files = append(files, file)

	}
	if files == nil {
		return make([]File, 0), nil
	}
	return &files, nil
}

func InsertMetadata(vaultID string, file File) (int, error) {
	// If created & modified are 0, set them to current time
	if file.Created == 0 {
		file.Created = time.Now().UnixMilli()
	}
	if file.Modified == 0 {
		file.Modified = time.Now().UnixMilli()
	}
	// Set previous files with the same path to not be newest
	_, err := db.Exec("UPDATE files SET newest = 0 WHERE path = ? AND newest = 1", file.Path)
	if err != nil {
		return 0, err
	}
	result, err := db.Exec(`INSERT INTO files (
		vault_id, path, hash, extension, size, created, modified, folder, deleted) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		vaultID, file.Path, file.Hash, file.Extension, file.Size, file.Created,
		file.Modified, file.Folder, file.Deleted)
	if err != nil {
		return 0, err
	}
	lastInsertID, err := result.LastInsertId()

	return int(lastInsertID), err
}

func GetFileData(uid int) (*[]byte, error) {
	var file []byte
	err := db.QueryRow("SELECT data FROM files WHERE uid = ?", uid).Scan(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func InsertData(uid int, data *[]byte) error {
	_, err := db.Exec("UPDATE files SET data = ? WHERE uid = ?", data, uid)
	return err
}

func DeleteVaultFile(path string) error {
	// Update all files with the same path to be deleted
	_, err := db.Exec("UPDATE files SET deleted = 1 WHERE path = ?", path)
	return err
}
