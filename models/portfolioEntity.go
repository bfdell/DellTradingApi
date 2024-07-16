package models

import (
	"time"

	"gorm.io/gorm"
)

type PortfolioEntity struct {
	UserID   uint      `gorm:"primaryKey"`
	Ticker   string    `gorm:"primaryKey"`
	AsOfTime time.Time `gorm:"primaryKey"`
	Shares   uint

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
