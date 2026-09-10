package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"caesar-quiz/internal/ws"
)

// WSHandler menangani koneksi WebSocket untuk game tertentu.
// Origin dibatasi sesuai allowedOrigins (bukan bebas semua).
func WSHandler(hub *ws.Hub, allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		gameIDStr := c.Param("gameID")
		gameID, err := uuid.Parse(gameIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gameID"})
			return
		}

		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return originAllowed(r, allowedOrigins)
			},
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

func originAllowed(r *http.Request, allowed []string) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true // non-browser client (CLI/tes)
	}
	for _, a := range allowed {
		if a == "*" || a == origin {
			return true
		}
	}
	return false
}
