package models

import "time"

// Модель для хранения Refresh токенов
type RefreshToken struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      string    `gorm:"not null"`
	TokenHash   string    `gorm:"not null;unique"`
	AccessToken string    `gorm:"not null;unique"`
	Blocked     bool      `gorm:"NOT NULL DEFAULT false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	ExpiresAt   time.Time `gorm:"not null"`
	IPAdress    string    `gorm:"not null"`
}
