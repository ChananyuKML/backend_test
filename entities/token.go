package entities

import "time"

type RefreshToken struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Token     string `gorm:"unique"`
	Revoked   bool
	ExpiresAt time.Time
}
