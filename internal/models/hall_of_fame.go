package models

import (
	"time"

	"github.com/google/uuid"
)

type HallOfFame struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PlayerID  uuid.UUID `gorm:"type:uuid;not null"`
	GameID    uuid.UUID `gorm:"type:uuid;not null"`
	TimeMs    int64     `gorm:"not null;default:0"` // waktu tercepat dalam ms
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Preload
	Player Player `gorm:"foreignKey:PlayerID"`
}
