package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"caesar-quiz/internal/models"
	"caesar-quiz/internal/services"
	"caesar-quiz/internal/ws"
)

// RegisterRoutes memasang seluruh endpoint. allowedOrigins dipakai WebSocket.
func RegisterRoutes(r *gin.Engine, db *gorm.DB, hub *ws.Hub, allowedOrigins []string) {
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Broadcast hint endpoint (khusus host)
	r.POST("/game/:gameID/hint", requireHost(db), func(c *gin.Context) {
		gameID := c.Param("gameID")
		var question models.Question
		if err := db.Where("game_id = ?", gameID).Order("created_at ASC").First(&question).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "question not found"})
			return
		}
		// Generate random Caesar mapping hint
		var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		shift := ((question.Shift % 26) + 26) % 26
		idx := int(time.Now().UnixNano() % int64(len(letters)))
		hint := string(letters[idx]) + " = " + string(letters[(idx+shift)%26])
		hub.BroadcastHint(question.GameID, hint)
		c.JSON(http.StatusOK, gin.H{"hint": hint})
	})

	// List players in a game
	r.GET("/game/:gameID/players", func(c *gin.Context) {
		gameID := c.Param("gameID")
		var players []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		if err := db.Raw("SELECT id, name FROM players WHERE game_id = ?", gameID).Scan(&players).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch players"})
			return
		}
		c.JSON(http.StatusOK, players)
	})

	// Daftar soal (tanpa plaintext) untuk halaman game.
	r.GET("/game/:gameID/questions", func(c *gin.Context) {
		gameID := c.Param("gameID")
		var questions []models.Question
		if err := db.Where("game_id = ?", gameID).Order("created_at ASC").Find(&questions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch questions"})
			return
		}
		resp := make([]gin.H, 0, len(questions))
		for _, q := range questions {
			resp = append(resp, gin.H{
				"id":     q.ID,
				"cipher": q.Cipher,
			})
		}
		c.JSON(http.StatusOK, resp)
	})

	r.POST("/game/:gameID/start", requireHost(db), StartGame(db, hub))
	r.POST("/game/:gameID/answer", SubmitAnswer(db, hub))
	r.POST("/game", CreateGame(db, hub))

	join := joinGameHandler(db)
	r.POST("/game/join", join)     // alias
	r.POST("/players/join", join)  // endpoint lama

	RegisterHallOfFameRoutes(r, db)

	// WebSocket
	r.GET("/ws/:gameID", WSHandler(hub, allowedOrigins))
}

// requireHost memastikan request membawa X-Host-Token yang cocok dengan game.
// Game lama tanpa HostToken tetap diizinkan (kompatibilitas).
func requireHost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		gameID := c.Param("gameID")
		var game models.Game
		if err := db.First(&game, "id = ?", gameID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
			c.Abort()
			return
		}
		if game.HostToken != "" {
			token := c.GetHeader("X-Host-Token")
			if token == "" {
				token = c.Query("host_token")
			}
			if token != game.HostToken {
				c.JSON(http.StatusForbidden, gin.H{"error": "invalid host token"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func joinGameHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Code string `json:"code"`
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		if !services.ValidateName(req.Name) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
			return
		}

		var game models.Game
		if err := db.First(&game, "code = ? AND is_active = true", req.Code).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
			return
		}

		player := models.Player{
			ID:     uuid.New(),
			GameID: game.ID,
			Name:   req.Name,
		}
		if err := db.Create(&player).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to join game"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"player_id": player.ID,
			"game_id":   game.ID,
		})
	}
}
