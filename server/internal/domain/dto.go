package domain

import (
	"mime/multipart"
	"time"
)

type InputUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}

type UpdateProfileDTO struct {
	UserId          uint
	CurrentPassword string
	NewPassword     string
	NewName         string
	FileHeader      *multipart.FileHeader
}

type UpdateUserDTO struct {
	UserId     uint
	NewName    string
	FileHeader *multipart.FileHeader
}

type InputRoomDTO struct {
	Name string `json:"name"`
}

type OutputMessageDTO struct {
	UserID     uint   `json:"user_id"`
	UserName   string `json:"user_name"`
	UserAvatar string `json:"user_avatar"`
	Content    string `json:"content"`
}

type InputMessageDTO struct {
	RoomID  uint   `json:"room_id"`
	UserID  uint   `json:"user_id"`
	Content string `json:"content"`
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type SessionDTO struct {
	ID           uint      `json:"id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserID       uint      `json:"user_id"`
	UserRole     string    `json:"user_role"`
}
