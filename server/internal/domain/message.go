package domain

import "time"

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RoomID    uint      `json:"room_id"`
	UserID    uint      `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
