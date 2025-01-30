package domain

import "time"

type Session struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserID       uint      `json:"user_id"`
}
