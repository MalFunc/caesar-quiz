package services

import (
	"time"

	"gorm.io/gorm"

	"caesar-quiz/internal/models"
	"github.com/google/uuid"
)

// UpdateHallOfFame menyimpan waktu tercepat ke leaderboard
func UpdateHallOfFame(db *gorm.DB, playerID, gameID uuid.UUID, timeMs int64) error {
	var hof models.HallOfFame
	err := db.First(&hof, "player_id = ? AND game_id = ?", playerID, gameID).Error

	if err == gorm.ErrRecordNotFound {
		// belum ada -> buat baru
		hof = models.HallOfFame{
			ID:        uuid.New(),
			PlayerID:  playerID,
			GameID:    gameID,
			TimeMs:    timeMs,
			CreatedAt: time.Now(),
		}
		return db.Create(&hof).Error
	} else if err != nil {
		return err
	}

	// update hanya jika waktu lebih cepat
	if timeMs < hof.TimeMs {
		hof.TimeMs = timeMs
		return db.Save(&hof).Error
	}

	return nil
}
