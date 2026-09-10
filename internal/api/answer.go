package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"caesar-quiz/internal/models"
	"caesar-quiz/internal/services"
	"caesar-quiz/internal/ws"
)

func SubmitAnswer(db *gorm.DB, hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			PlayerID   string `json:"player_id"`
			QuestionID string `json:"question_id"`
			Answer     string `json:"answer"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		playerID, err := uuid.Parse(req.PlayerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player_id"})
			return
		}
		questionID, err := uuid.Parse(req.QuestionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question_id"})
			return
		}

		var question models.Question
		if err := db.First(&question, "id = ?", questionID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "question not found"})
			return
		}

		// Pastikan pertanyaan ini milik game dari si pemain (cegah cross-game).
		var player models.Player
		if err := db.First(&player, "id = ?", playerID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
			return
		}
		if player.GameID != question.GameID {
			c.JSON(http.StatusForbidden, gin.H{"error": "player does not belong to this game"})
			return
		}

		// Game harus aktif & sudah dimulai.
		var game models.Game
		if err := db.First(&game, "id = ?", question.GameID).Error; err != nil || !game.IsActive {
			c.JSON(http.StatusForbidden, gin.H{"error": "game is not active"})
			return
		}
		if question.StartedAt.IsZero() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "game has not started yet"})
			return
		}

		// Cegah submit ulang / brute force: kalau sudah pernah benar, balas yang lama.
		var existing models.Answer
		err = db.Where("player_id = ? AND question_id = ? AND is_correct = ?", playerID, questionID, true).
			First(&existing).Error
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"correct": true, "time_ms": existing.TimeMs, "duplicate": true})
			return
		} else if err != gorm.ErrRecordNotFound {
			log.Printf("answer: lookup existing failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit answer"})
			return
		}

		isCorrect := question.Plaintext == req.Answer
		timeMs := time.Since(question.StartedAt).Milliseconds()
		if timeMs < 0 {
			timeMs = 0
		}

		answer := models.Answer{
			ID:         uuid.New(),
			PlayerID:   playerID,
			QuestionID: questionID,
			Answer:     req.Answer,
			IsCorrect:  isCorrect,
			AnsweredAt: time.Now(),
			TimeMs:     timeMs,
		}

		if err := db.Create(&answer).Error; err != nil {
			log.Printf("answer: create failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit answer"})
			return
		}

		// update leaderboard
		if isCorrect {
			if err := services.UpdateHallOfFame(db, playerID, question.GameID, answer.TimeMs); err != nil {
				log.Printf("answer: update hall of fame failed: %v", err)
			}
			hub.BroadcastAnswer(question.GameID)
		}

		c.JSON(http.StatusOK, gin.H{"correct": isCorrect, "time_ms": answer.TimeMs})
	}
}
