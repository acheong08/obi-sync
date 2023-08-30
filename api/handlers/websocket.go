package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/acheong08/obsidian-sync/database"
	"github.com/acheong08/obsidian-sync/utilities"
	"github.com/acheong08/obsidian-sync/vault"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// ReadBufferSize:  1024,
	// WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// A channel manager: push messages to all clients on the same vault
type channelManager struct {
	clients map[*websocket.Conn]bool
}

func (cm *channelManager) AddClient(ws *websocket.Conn) {
	cm.clients[ws] = true
}

func (cm *channelManager) RemoveClient(ws *websocket.Conn) {
	delete(cm.clients, ws)
}

func (cm *channelManager) IsEmpty() bool {
	return len(cm.clients) == 0
}

func (cm *channelManager) Broadcast(msg any) {
	for client := range cm.clients {
		client.WriteJSON(msg)
	}
}

// map[vaultID]ChannelManager
var channels map[string]*channelManager = make(map[string]*channelManager)

func getMsg(ws *websocket.Conn) ([]byte, error) {
	msgType, msg, err := ws.ReadMessage()
	if err != nil {
		return nil, err
	}
	if msgType != websocket.TextMessage {
		return nil, errors.New("message type must be text")
	}
	return msg, nil
}

func WsHandler(c *gin.Context) {
	// Upgrade protocol to websocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	// Receive initialization message
	msg, err := getMsg(ws)
	if err != nil {
		// Send error message
		ws.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	// Parse initialization message as JSON
	connectionInfo, connectedVault, err := initHandler(msg, c.MustGet("db").(*database.Database))
	if err != nil {
		ws.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	// No errors: ok response
	ws.WriteJSON(gin.H{"res": "ok"})
	var version int = utilities.ToInt(connectionInfo.Version)
	if err != nil {
		ws.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	if connectedVault.Version > version {
		vaultFiles, err := vault.GetVaultFiles(connectedVault.ID)
		if err != nil {
			ws.WriteJSON(gin.H{"error": err.Error()})
			return
		}
		for _, file := range *vaultFiles {
			ws.WriteJSON(gin.H{
				"op": "push", "path": file.Path,
				"hash": file.Hash, "size": file.Size,
				"ctime": file.Created, "mtime": file.Modified, "folder": file.Folder,
				"deleted": file.Deleted, "device": "insignificantv5", "uid": file.UID})

		}

	}
	var version_bumped bool = false
	ws.WriteJSON(gin.H{"op": "ready", "version": connectedVault.Version})

	defer vault.Snapshot(connectedVault.ID)

	dbConnection := c.MustGet("db").(*database.Database)
	if connectedVault.Version < version {
		dbConnection.SetVaultVersion(connectedVault.ID, version)
	}

	// Check if vaultID is in channels
	if _, ok := channels[connectedVault.ID]; !ok {
		// Create new channel manager
		channels[connectedVault.ID] = &channelManager{
			clients: make(map[*websocket.Conn]bool),
		}
	}
	channels[connectedVault.ID].AddClient(ws)
	defer channels[connectedVault.ID].RemoveClient(ws)
	defer func() {
		if channels[connectedVault.ID].IsEmpty() {
			delete(channels, connectedVault.ID)
		}
	}()

	// Inifinite loop to handle messages
	type message struct {
		Op string `json:"op" binding:"required"` // Operation
	}
	for {
		msg, err := getMsg(ws)
		if err != nil {
			// Send error message
			ws.WriteJSON(gin.H{"error": err.Error()})
			return
		}
		var m message
		err = json.Unmarshal(msg, &m)
		if err != nil {
			ws.WriteJSON(gin.H{"error": err.Error()})
			return
		}
		switch m.Op {
		case "size":
			size, err := vault.GetVaultSize(connectedVault.ID)
			if err != nil {
				ws.WriteJSON(gin.H{"error": err.Error()})
				return
			}
			ws.WriteJSON(gin.H{
				"res":   "ok",
				"size":  size,
				"limit": 10737418240,
			})
		case "pull":
			var pull struct {
				UID any `json:"uid" binding:"required"`
			}
			err = json.Unmarshal(msg, &pull)
			if err != nil {
				ws.WriteJSON(gin.H{"error": err.Error()})
				return
			}
			var uid int = utilities.ToInt(pull.UID)
			file, err := vault.GetFile(uid)
			if err != nil {
				ws.WriteJSON(gin.H{"error": err.Error()})
				return
			}
			var pieces int8 = 0
			if file.Size != 0 {
				pieces = 1
			}
			ws.WriteJSON(gin.H{
				"hash": file.Hash, "size": file.Size, "pieces": pieces,
			})
			if file.Size != 0 {
				ws.WriteMessage(websocket.BinaryMessage, file.Data)
			}
		case "push":
			var metadata struct {
				UID          int    `json:"uid"`
				Op           string `json:"op" binding:"required"` // Operation
				Path         string `json:"path" binding:"required"`
				Extension    string `json:"extension" binding:"required"`
				Hash         string `json:"hash" binding:"required"`
				CreationTime int64  `json:"ctime" binding:"required"`
				ModifiedTime int64  `json:"mtime" binding:"required"`
				Folder       bool   `json:"folder" binding:"required"`
				Deleted      bool   `json:"deleted" binding:"required"`
				Size         int    `json:"size,omitempty" binding:"required"`
				Pieces       int    `json:"pieces,omitempty" binding:"required"`
			}
			err = json.Unmarshal(msg, &metadata)
			if err != nil {
				ws.WriteJSON(gin.H{"error": err.Error()})
				return
			}
			if metadata.Deleted {
				err := vault.DeleteVaultFile(metadata.Path)
				if err != nil {
					ws.WriteJSON(gin.H{"error": err.Error()})
					return
				}
			}
			vaultUID, err := vault.InsertMetadata(connectedVault.ID, vault.File{
				Path:      metadata.Path,
				Hash:      metadata.Hash,
				Extension: metadata.Extension,
				Size:      int64(metadata.Size),
				Created:   metadata.CreationTime,
				Modified:  metadata.ModifiedTime,
				Folder:    metadata.Folder,
				Deleted:   metadata.Deleted,
			})
			if err != nil {
				ws.WriteJSON(gin.H{"error": err.Error()})
				return
			}
			var fullBinary []byte
			for i := 0; i < metadata.Pieces; i++ {
				ws.WriteJSON(gin.H{"res": "next"})
				// Read bytes
				msgType, msg, err := ws.ReadMessage()
				if err != nil {
					ws.WriteJSON(gin.H{"error": err.Error()})
					return
				}
				if msgType != websocket.BinaryMessage {
					ws.WriteJSON(gin.H{"error": "message type must be binary"})
					return
				}
				fullBinary = append(fullBinary, msg...)
			}
			err = vault.InsertData(vaultUID, &fullBinary)
			if err != nil {
				ws.WriteJSON(gin.H{"error": err.Error()})
				return
			}
			metadata.UID = int(vaultUID)
			// Broadcast to all clients
			channels[connectedVault.ID].Broadcast(metadata)
			if !version_bumped {
				dbConnection.SetVaultVersion(connectedVault.ID, version+1)
				version_bumped = true
			}
			ws.WriteJSON(gin.H{"op": "ok"})
		case "history":
			var history struct {
				Last any    `json:"last"` // I have no idea what this is for
				Path string `json:"path" binding:"required"`
			}
			err = json.Unmarshal(msg, &history)
			if err != nil {
				ws.WriteJSON(gin.H{"error": err.Error()})
				return
			}
			files, err := vault.GetFileHistory(history.Path)
			if err != nil {
				ws.WriteJSON(gin.H{"error": err.Error()})
				return
			}
			ws.WriteJSON(gin.H{"items": files, "more": false})

		case "ping":
			ws.WriteJSON(gin.H{"op": "pong"})
		case "deleted":
			files, err := vault.GetDeletedFiles()
			if err != nil {
				ws.WriteJSON(gin.H{"error": err.Error()})
				return
			}
			ws.WriteJSON(gin.H{"items": files})
		case "restore":
			var restore struct {
				UID any `json:"uid" binding:"required"`
			}
			err = json.Unmarshal(msg, &restore)
			if err != nil {
				ws.WriteJSON(gin.H{"error": err.Error()})
				return
			}
			var uid int = utilities.ToInt(restore.UID)
			err = json.Unmarshal(msg, &restore)
			if err != nil {
				ws.WriteJSON(gin.H{"error": err.Error()})
				return
			}
			file, err := vault.RestoreFile(uid)
			if err != nil {
				ws.WriteJSON(gin.H{"error": err.Error()})
				return
			}
			file.Op = "push"
			channels[connectedVault.ID].Broadcast(file)
			ws.WriteJSON(gin.H{"res": "ok"})
		default:
			log.Println("Unknown operation:", m.Op)
			log.Println("Data: ", string(msg))
		}

	}

}

type initializationRequest struct {
	Op      string `json:"op" binding:"required"` // Operation
	Token   string `json:"token" binding:"required"`
	Id      string `json:"id" binding:"required"`      // Vault ID
	KeyHash string `json:"keyhash" binding:"required"` // Hash of password & salt
	Version any    `json:"version" binding:"required"`
	Initial bool   `json:"initial" binding:"required"`
	Device  string `json:"device" binding:"required"`
}

func initHandler(req []byte, dbConnection *database.Database) (*initializationRequest, *vault.Vault, error) {

	var initial initializationRequest
	err := json.Unmarshal(req, &initial)
	if err != nil {
		return nil, nil, err
	}

	// Validate token and key hash
	email, err := utilities.GetJwtEmail(initial.Token)
	if err != nil {
		return nil, nil, err
	}

	vault, err := dbConnection.GetVault(initial.Id, initial.KeyHash)
	if err != nil {
		return nil, nil, err
	}
	if !dbConnection.HasAccessToVault(vault.ID, email) {
		return nil, nil, fmt.Errorf("no access to vault")
	}
	return &initial, vault, nil
}
