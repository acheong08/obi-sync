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
		CREATE TABLE IF NOT EXISTS file_metadata (
			uid INTEGER PRIMARY KEY,
			vault_id TEXT,
			path TEXT,
			hash TEXT,
			size INTEGER,
			created INTEGER,
			modified INTEGER,
			folder INTEGER,
			deleted INTEGER,
		);
		
		CREATE TABLE IF NOT EXISTS file (
			path TEXT PRIMARY KEY,
			data BLOB
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func GetVaultFiles(vaultID string) (*[]FileMetadata, error) {
	rows, err := db.Query("SELECT * FROM file_metadata WHERE vault_id = ?", vaultID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var files []FileMetadata
	for rows.Next() {
		var file FileMetadata
		err = rows.Scan(&file.UID, &file.Path, &file.Hash, &file.Size, &file.Created, &file.Modified, &file.Folder, &file.Deleted)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return &files, nil
}

func GetVaultFile(uid int) (*FileMetadata, error) {
	var file FileMetadata
	err := db.QueryRow("SELECT * FROM file_metadata WHERE uid = ?", uid).Scan(&file.UID, &file.Path, &file.Hash, &file.Size, &file.Created, &file.Modified, &file.Folder, &file.Deleted)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func InsertVaultFile(vaultID string, file FileMetadata) error {
	_, err := db.Exec("INSERT INTO file_metadata (vault_id, path, hash, size, created, modified, folder, deleted) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", vaultID, file.Path, file.Hash, file.Size, file.Created, file.Modified, file.Folder, file.Deleted)
	return err
}

func GetFile(path string) (*[]byte, error) {
	var file []byte
	err := db.QueryRow("SELECT data FROM file WHERE path = ?", path).Scan(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func PushFile(path string, data *[]byte) error {
	// Overwrite file if it already exists
	_, err := db.Exec("INSERT INTO file (path, data) VALUES (?, ?) ON CONFLICT(path) DO UPDATE SET data = ?", path, data, data)
	return err
}
