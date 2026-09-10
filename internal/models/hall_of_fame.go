package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HallOfFame struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	PlayerID  uuid.UUID `gorm:"type:uuid;not null"`
	GameID    uuid.UUID `gorm:"type:uuid;not null"`
	TimeMs    int64     `gorm:"not null;default:0"` // waktu tercepat dalam ms
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Preload
	Player Player `gorm:"foreignKey:PlayerID"`
}

func (h *HallOfFame) BeforeCreate(tx *gorm.DB) (err error) {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return
}
