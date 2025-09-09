package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Question struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	GameID    uuid.UUID `gorm:"type:uuid;not null"`
	Plaintext string    `gorm:"not null"`
	Cipher    string    `gorm:"not null"`
	Shift     int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	StartedAt time.Time 
	Game    Game     `gorm:"foreignKey:GameID"`
	Answers []Answer `gorm:"foreignKey:QuestionID"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) (err error) {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return
}
