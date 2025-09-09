package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"caesar-quiz/internal/models"
	"github.com/google/uuid"
	"caesar-quiz/internal/ws"
	"time"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, hub *ws.Hub) {
	// Broadcast hint endpoint
	r.POST("/game/:gameID/hint", func(c *gin.Context) {
		gameID := c.Param("gameID")
		var question models.Question
		if err := db.Where("game_id = ?", gameID).Order("created_at ASC").First(&question).Error; err != nil {
			c.JSON(404, gin.H{"error": "question not found"})
			return
		}
		// Generate random Caesar mapping hint
		var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		shift := question.Shift
		idx := int(time.Now().UnixNano() % int64(len(letters)))
		plain := letters[idx]
		cipher := letters[(idx+shift)%26]
		hint := string(plain) + " = " + string(cipher)
		hub.BroadcastHint(question.GameID, hint)
		c.JSON(200, gin.H{"hint": hint})
	})
	// List players in a game
	r.GET("/game/:gameID/players", func(c *gin.Context) {
		gameID := c.Param("gameID")
		var players []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		if err := db.Raw("SELECT id, name FROM players WHERE game_id = ?", gameID).Scan(&players).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to fetch players"})
			return
		}
		c.JSON(200, players)
	})
	r.POST("/game/:gameID/start", StartGame(db, hub))
	RegisterPlayerRoutes(r, db)
	r.POST("/game/:gameID/answer", SubmitAnswer(db, hub))
	RegisterHallOfFameRoutes(r, db)
	r.POST("/game", CreateGame(db, hub))
	// Add /game/join endpoint for joining a game (alias for /players/join)
	r.POST("/game/join", func(c *gin.Context) {
		var req struct {
			Code string `json:"code"`
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}
		if len(req.Name) == 0 || len(req.Name) > 255 {
			c.JSON(400, gin.H{"error": "invalid name"})
			return
		}
		var game models.Game
		if err := db.First(&game, "code = ? AND is_active = true", req.Code).Error; err != nil {
			c.JSON(404, gin.H{"error": "game not found"})
			return
		}
		player := models.Player{
			ID:     uuid.New(),
			GameID: game.ID,
			Name:   req.Name,
		}
		if err := db.Create(&player).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to join game"})
			return
		}
		c.JSON(200, gin.H{
			"player_id": player.ID,
			"game_id":   game.ID,
		})
	})
	// WebSocket
	r.GET("/ws/:gameID", WSHandler(hub))
}
