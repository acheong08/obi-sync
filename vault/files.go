package vault

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sync"
)

func init() {
	// Make data directory if it doesn't exist
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}
	path := "data/files_metadata.gob"
	file, err := os.Open(path)
	if err != nil {
		log.Println("Could not open file", path) // A new file will be made
		return
	}
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&files_metadata)
	if err != nil {
		log.Println("Could not decode file", path)
		return
	}
}

type filesMetadataMutex struct {
	FilesMetadata map[string][]FileMetadata
	Mutex         sync.Mutex
}

var files_metadata = &filesMetadataMutex{}

func GetVaultFiles(vaultID string) []FileMetadata {
	files_metadata.Mutex.Lock()
	defer files_metadata.Mutex.Unlock()
	return files_metadata.FilesMetadata[vaultID]
}

func InsertVaultFile(vaultID string, file FileMetadata) {
	files_metadata.Mutex.Lock()
	defer files_metadata.Mutex.Unlock()
	files_metadata.FilesMetadata[vaultID] = append(files_metadata.FilesMetadata[vaultID], file)
}

func WriteVaultFiles() error {
	// Write files_metadata to gob
	path := "data/files_metadata.gob"
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(files_metadata)
	return err
}

func GetFile(path string) (*File, error) {
	path = fmt.Sprintf("data/%s.gob", path)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	data := &File{}
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(data)
	return data, err
}

func PushFile(path string, data []byte) error {
	path = fmt.Sprintf("data/%s.gob", path)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(&File{Data: data})
	return err
}
