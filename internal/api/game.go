
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"time"

	"caesar-quiz/internal/models"
	"caesar-quiz/internal/services"
	"caesar-quiz/internal/ws"
)

// StartGame memulai game dan broadcast soal pertama ke peserta
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
		db.Save(&question) 
		 // Broadcast soal ke semua peserta via WebSocket
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			   CreatedAt: time.Now(),
		   }

		   if err := db.Create(&game).Error; err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
		   }

		   // Insert question if target is provided
		   if req.Target != "" {
			   cipher := services.CaesarEncrypt(req.Target, game.Shift)
			   question := models.Question{
				   ID:        uuid.New(),
				   GameID:    game.ID,
				   Plaintext: req.Target,
				   Cipher:    cipher,
				   Shift:     game.Shift,
				   CreatedAt: time.Now(),
			   }
			   if err := db.Create(&question).Error; err != nil {
				   c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create question"})
				   return
			   }
		   }

		   c.JSON(http.StatusOK, gin.H{
			   "game_id": game.ID,
			   "code":    game.Code,
			   "shift":   game.Shift,
		   })
	}
}
