package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Player struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	GameID    uuid.UUID `gorm:"type:uuid;not null"`
	Name      string    `gorm:"size:255;not null"` // only huruf, validasi di layer API
	JoinedAt  time.Time `gorm:"autoCreateTime"`

	Game    Game     `gorm:"foreignKey:GameID"`
	Answers []Answer `gorm:"foreignKey:PlayerID"`
}

func (p *Player) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
