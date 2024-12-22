package model

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex" json:"email"`
	PasswordHash string    `json:"password_hash"`
	Name         string    `json:"name"`
	PhotoURL     string    `json:"photo_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Room struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Message struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RoomID     uint      `json:"room_id"`
	UserID     uint      `json:"user_id"`
	UserName   string    `json:"user_name"`
	UserAvatar string    `json:"user_avatar"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
