package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Answer struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	PlayerID   uuid.UUID `gorm:"type:uuid;not null"`
	QuestionID uuid.UUID `gorm:"type:uuid;not null"`
	Answer     string    `gorm:"not null"`
	IsCorrect  bool      `gorm:"not null;default:false"`
	AnsweredAt time.Time `gorm:"autoCreateTime"`
	TimeMs     int64     // waktu player menjawab (ms)

	// Relations
	Player   Player   `gorm:"foreignKey:PlayerID"`
	Question Question `gorm:"foreignKey:QuestionID"`
}

func (a *Answer) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
