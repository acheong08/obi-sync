package handlers

import (
	"encoding/json"
	"errors"

	"github.com/acheong08/obsidian-sync/database"
	"github.com/acheong08/obsidian-sync/utilities"
	"github.com/acheong08/obsidian-sync/vault"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

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
		c.JSON(500, gin.H{
			"message": "error upgrading to websocket",
		})
		return
	}
	defer ws.Close()

	// Recieve initialization message
	msg, err := getMsg(ws)
	if err != nil {
		// Send error message
		ws.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	// Parse initialization message as JSON
	connectionInfo, _, connectedVault, err := initHandler(msg, c.MustGet("db").(*database.Database))
	if err != nil {
		ws.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	// No errors: ok response
	ws.WriteJSON(gin.H{"res": "ok"})
	if connectionInfo.Initial {
		if connectionInfo.Version != connectedVault.Version {
			vaultFiles, err := vault.GetVaultFiles(connectedVault.ID)
			if err != nil {
				ws.WriteJSON(gin.H{"error": err.Error()})
				return
			}
			for _, file := range vaultFiles {
				if !file.Deleted {
					ws.WriteJSON(gin.H{
						"op": "push", "path": file.Path,
						"hash": file.Hash, "size": file.Size,
						"ctime": file.Created, "mtime": file.Modified, "folder": file.Folder,
						"deleted": file.Deleted, "device": "insignificantv5", "uid": file.UID})
				}
			}
		}
	}
	ws.WriteJSON(gin.H{"op": "ready", "version": connectedVault.Version})

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
			ws.WriteJSON(gin.H{
				"res":   "ok",
				"size":  connectedVault.Size,
				"limit": 10737418240,
			})
		}
	}

}

type initializationRequest struct {
	Op      string `json:"op" binding:"required"` // Operation
	Token   string `json:"token" binding:"required"`
	Id      string `json:"id" binding:"required"`      // Vault ID
	KeyHash string `json:"keyhash" binding:"required"` // Hash of password & salt
	Version int    `json:"version" binding:"required"`
	Initial bool   `json:"initial" binding:"required"`
	Device  string `json:"device" binding:"required"`
}

func initHandler(req []byte, dbConnection *database.Database) (*initializationRequest, string, *vault.Vault, error) {

	var initial initializationRequest
	err := json.Unmarshal(req, &initial)
	if err != nil {
		return nil, "", nil, err
	}

	// Validate token and key hash
	email, err := utilities.GetJwtEmail(initial.Token)
	if err != nil {
		return nil, "", nil, err
	}

	vault, err := dbConnection.GetVault(initial.Id, initial.KeyHash)
	if err != nil {
		return nil, "", nil, err
	}
	return &initial, email, vault, nil
}