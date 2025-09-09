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

func SubmitAnswer(db *gorm.DB, hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			PlayerID   string `json:"player_id"`
			QuestionID string `json:"question_id"`
			Answer     string `json:"answer"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		playerID, _ := uuid.Parse(req.PlayerID)
		questionID, _ := uuid.Parse(req.QuestionID)

		var question models.Question
		if err := db.First(&question, "id = ?", questionID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "question not found"})
			return
		}

		isCorrect := question.Plaintext == req.Answer
		timeMs := int64(time.Since(question.StartedAt) / time.Second)
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// update leaderboard
		if isCorrect {
			services.UpdateHallOfFame(db, playerID, question.GameID, answer.TimeMs)
			hub.BroadcastAnswer(question.GameID) // <- Hanya gameID
		}

		c.JSON(http.StatusOK, gin.H{"correct": isCorrect})
	}
}
