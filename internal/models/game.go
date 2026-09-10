package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Game struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Code      string    `gorm:"uniqueIndex;not null;size:16"` // kode unik game
	Shift     int       `gorm:"not null"`                     // Caesar shift (bisa random/custom)
	HostToken string    `gorm:"size:64"`                      // token rahasia host (start/hint)
	IsActive  bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Players   []Player   `gorm:"foreignKey:GameID"`
	Questions []Question `gorm:"foreignKey:GameID"`
}

func (g *Game) BeforeCreate(tx *gorm.DB) (err error) {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return
}
