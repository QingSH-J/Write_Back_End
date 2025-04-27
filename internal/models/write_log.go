package models

import (
	"time"

	"github.com/google/uuid"
)

type WriteLog struct {
	ID         uuid.UUID `gorm:"primary_key" json:"id"`
	UserID     uuid.UUID `gorm:"index;not null" json:"userId"`
	Date       time.Time `gorm:"type:date;not null" json:"date"` //remove the uniqueIndex to support multiple logs for the same day
	WordsCount int       `gorm:"not null" json:"wordsCount"`
	Content    string    `gorm:"type:text;not null" json:"content,omitempty"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
