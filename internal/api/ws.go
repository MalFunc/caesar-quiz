package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"caesar-quiz/internal/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WSHandler menangani koneksi WebSocket untuk game tertentu
func WSHandler(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		gameIDStr := c.Param("gameID")
		gameID, err := uuid.Parse(gameIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gameID"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade websocket"})
			return
		}

		hub.RegisterClient(gameID, conn)

		// tunggu sampai koneksi ditutup
		for {
			if _, _, err := conn.NextReader(); err != nil {
				hub.UnregisterClient(gameID, conn)
				break
			}
		}
	}
}
