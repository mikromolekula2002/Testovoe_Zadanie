package models

import "time"

// Модель для хранения Refresh токенов
type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"not null"`
	TokenHash string    `gorm:"not null;unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	IPAdress  string    `gorm:"not null"`
}
