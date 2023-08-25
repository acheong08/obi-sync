package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
)

func ListSubscriptions(c *gin.Context) {
	c.JSON(200, gin.H{
		"business": nil,
		"publish":  nil,
		"sync": gin.H{
			"earlybird": false,
			// Always expires in 1 year
			"expiry_ts": time.Now().UnixMilli() + time.Hour.Milliseconds()*24*365,
			"renew":     "",
		},
	})
}
