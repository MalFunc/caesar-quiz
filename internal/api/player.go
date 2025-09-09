package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"caesar-quiz/internal/models"
)

func RegisterPlayerRoutes(r *gin.Engine, db *gorm.DB) {
	players := r.Group("/players")
	{
		players.POST("/join", func(c *gin.Context) {
			var req struct {
				Code string `json:"code"`
				Name string `json:"name"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}

			if len(req.Name) == 0 || len(req.Name) > 255 {
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
		})
	}
}
