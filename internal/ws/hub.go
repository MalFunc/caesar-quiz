
package ws

import (
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"caesar-quiz/internal/models"
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients   map[uuid.UUID]map[*websocket.Conn]bool // gameID -> connections
	register   chan subscription
	unregister chan subscription
	broadcast  chan broadcastMessage
	db         *gorm.DB
	mu         sync.Mutex
}

type subscription struct {
	gameID uuid.UUID
	conn   *websocket.Conn
}

type broadcastMessage struct {
	gameID uuid.UUID
	data   interface{}
}

func NewHub(db *gorm.DB) *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]map[*websocket.Conn]bool),
		register:   make(chan subscription),
		unregister: make(chan subscription),
		broadcast:  make(chan broadcastMessage),
		db:         db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case sub := <-h.register:
			h.mu.Lock()
			if h.clients[sub.gameID] == nil {
				h.clients[sub.gameID] = make(map[*websocket.Conn]bool)
			}
			h.clients[sub.gameID][sub.conn] = true
			h.mu.Unlock()

		case sub := <-h.unregister:
			h.mu.Lock()
			if conns, ok := h.clients[sub.gameID]; ok {
				if _, exists := conns[sub.conn]; exists {
					delete(conns, sub.conn)
					sub.conn.Close()
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.Lock()
			for conn := range h.clients[msg.gameID] {
				conn.WriteJSON(msg.data)
			}
			h.mu.Unlock()
		}
	}
}

// RegisterClient menambahkan koneksi ke game tertentu
func (h *Hub) RegisterClient(gameID uuid.UUID, conn *websocket.Conn) {
	h.register <- subscription{gameID, conn}
}

// UnregisterClient menghapus koneksi
func (h *Hub) UnregisterClient(gameID uuid.UUID, conn *websocket.Conn) {
	h.unregister <- subscription{gameID, conn}
}

// BroadcastQuestion broadcast soal ke semua player
func (h *Hub) BroadcastQuestion(gameID uuid.UUID, question models.Question) {
	   h.broadcast <- broadcastMessage{
		   gameID: gameID,
		   data: map[string]interface{}{
			   "type":      "question",
			   "id":        question.ID,
			   "cipher":    question.Cipher,
		   },
	   }
}

// BroadcastAnswer broadcast leaderboard terbaru
func (h *Hub) BroadcastAnswer(gameID uuid.UUID) {
	var leaderboard []models.HallOfFame
	h.db.Preload("Player").Where("game_id = ?", gameID).Order("time_ms ASC").Find(&leaderboard)

	h.broadcast <- broadcastMessage{
		gameID: gameID,
		data: map[string]interface{}{
			"type":        "leaderboard",
			"leaderboard": buildLeaderboardResponse(leaderboard),
		},
	}
}

func buildLeaderboardResponse(leaderboard []models.HallOfFame) []map[string]interface{} {
	resp := []map[string]interface{}{}
	for i, entry := range leaderboard {
		resp = append(resp, map[string]interface{}{
			"rank":      i + 1,
			"player":    entry.Player.Name,
			"time_ms":   entry.TimeMs,
			"player_id": entry.PlayerID,
		})
	}
	return resp
}
// BroadcastHint broadcast hint Caesar ke semua player
func (h *Hub) BroadcastHint(gameID uuid.UUID, hint string) {
	h.broadcast <- broadcastMessage{
		gameID: gameID,
		data: map[string]interface{}{
			"type": "hint",
			"hint": hint,
		},
	}
}
func (h *Hub) GetDB() *gorm.DB {
    return h.db
}
