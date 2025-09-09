package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Game struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Code      string    `gorm:"uniqueIndex;not null"` // kode unik game
	Shift     int       `gorm:"not null"`             // Caesar shift (bisa random/custom)
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
