package vault

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite", "vault.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create tables if they don't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			uid INTEGER PRIMARY KEY AUTOINCREMENT,
			vault_id TEXT,
			path TEXT,
			extension TEXT,
			size INTEGER,
			created INTEGER,
			modified INTEGER,
			folder INTEGER,
			deleted INTEGER,
			data BLOB
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
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
	rows, err := db.Query("SELECT uid, path, hash, extension, size, created, modified, folder, deleted FROM files WHERE vault_id = ? AND deleted = 0", vaultID)
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
	err := db.QueryRow("SELECT hash, size FROM files WHERE uid = ?", uid).Scan(&file.Hash, &file.Size)
	return &file, err
}

func GetFileHistory(path string) (*[]File, error) {
	rows, err := db.Query("SELECT uid, path, hash, extension, size, created, modified, folder, deleted FROM files WHERE path = ? AND deleted = 0", path)
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

func InsertMetadata(vaultID string, file File) (int, error) {
	result, err := db.Exec(`INSERT OR REPLACE INTO files (
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
	// Overwrite file if it already exists
	_, err := db.Exec("INSERT OR REPLACE INTO files (uid, data) VALUES (?, ?)", uid, *data)
	return err
}

// func DeleteVaultFile(uid int) error {
// 	_, err := db.Exec("DELETE FROM files WHERE uid = ?", uid)
// 	return err
// }
