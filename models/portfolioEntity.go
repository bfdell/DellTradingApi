package models

import (
	"time"

	"gorm.io/gorm"
)

type PortfolioEntity struct {
	UserID uint   `gorm:"primaryKey"`
	Ticker string `gorm:"primaryKey"`

	CreatedAt time.Time `gorm:"primaryKey"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Shares uint
	Cash   float64 `gorm:"type:numeric"`
}
