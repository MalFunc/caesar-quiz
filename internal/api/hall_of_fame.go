package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"caesar-quiz/internal/models"
)

func RegisterHallOfFameRoutes(r *gin.Engine, db *gorm.DB) {
	hof := r.Group("/hall-of-fame")
	{
		hof.GET("/:gameID", func(c *gin.Context) {
			gameID := c.Param("gameID")

			var leaderboard []models.HallOfFame
			if err := db.Preload("Player").Where("game_id = ?", gameID).Order("time_ms ASC").Find(&leaderboard).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch leaderboard"})
				return
			}

			// transform data supaya hanya kirim nama + waktu + rank
			response := []gin.H{}
			for i, entry := range leaderboard {
				response = append(response, gin.H{
					"rank":      i + 1,
					"player":    entry.Player.Name,
					"time_ms":   entry.TimeMs,
					"player_id": entry.PlayerID,
				})
			}

			c.JSON(http.StatusOK, response)
		})
	}
}
