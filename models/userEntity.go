package models

import (
	"gorm.io/gorm"
)

const JWT_EXPIRATION_DAYS uint32 = 2

type UserEntity struct {
	gorm.Model
	Email     string  `gorm:"uniqueIndex"`
	FirstName string  `gorm:"varchar(255);not null"`
	LastName  string  `gorm:"varchar(255);not null"`
	Password  string  `gorm:"not null"`
	Cash      float64 `gorm:"type:numeric;default:100000"`
}
