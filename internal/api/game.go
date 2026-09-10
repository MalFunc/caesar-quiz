package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"caesar-quiz/internal/models"
	"caesar-quiz/internal/services"
	"caesar-quiz/internal/ws"
)

const maxTargetLen = 200

// StartGame memulai game dan broadcast soal pertama ke peserta.
// Dilindungi host token (lihat requireHost di routes.go).
func StartGame(db *gorm.DB, hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		gameID := c.Param("gameID")
		var game models.Game
		if err := db.First(&game, "id = ?", gameID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
			return
		}

		// Ambil soal pertama untuk game ini
		var question models.Question
		if err := db.Where("game_id = ?", gameID).Order("created_at ASC").First(&question).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "no question found for this game"})
			return
		}
		question.StartedAt = time.Now()
		if err := db.Save(&question).Error; err != nil {
			log.Printf("start game: save question failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start game"})
			return
		}

		hub.BroadcastQuestion(game.ID, question)

		c.JSON(http.StatusOK, gin.H{"message": "game started", "question_id": question.ID})
	}
}

func CreateGame(db *gorm.DB, hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name   string `json:"name"`
			Shift  int    `json:"shift"`
			Target string `json:"target"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if req.Shift != 0 && !services.ValidShift(req.Shift) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "shift must be between 1 and 25"})
			return
		}
		target := strings.TrimSpace(req.Target)
		if len(target) > maxTargetLen {
			c.JSON(http.StatusBadRequest, gin.H{"error": "target too long"})
			return
		}

		shift := req.Shift
		if shift == 0 {
			shift = services.RandomShift()
		}

		game := models.Game{
			ID:        uuid.New(),
			Code:      services.GenerateGameCode(),
			Shift:     shift,
			HostToken: uuid.NewString(),
			CreatedAt: time.Now(),
		}

		if err := db.Create(&game).Error; err != nil {
			log.Printf("create game failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create game"})
			return
		}

		// Insert question if target is provided
		var questionID, cipher string
		if target != "" {
			cipher = services.CaesarEncrypt(target, game.Shift)
			question := models.Question{
				ID:        uuid.New(),
				GameID:    game.ID,
				Plaintext: target,
				Cipher:    cipher,
				Shift:     game.Shift,
				CreatedAt: time.Now(),
			}
			if err := db.Create(&question).Error; err != nil {
				log.Printf("create question failed: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create question"})
				return
			}
			questionID = question.ID.String()
		}

		c.JSON(http.StatusOK, gin.H{
			"game_id":     game.ID,
			"code":        game.Code,
			"shift":       game.Shift,
			"host_token":  game.HostToken,
			"question_id": questionID,
			"cipher":      cipher,
		})
	}
}
