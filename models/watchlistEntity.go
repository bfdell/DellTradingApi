package models

import (
	"time"

	"gorm.io/gorm"
)

type WatchlistEntity struct {
	UserID uint   `gorm:"primaryKey"`
	Ticker string `gorm:"primaryKey"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
