package repository

import "github.com/Nikita-Mihailuk/gochat/server/internal/domain"

type User interface {
	Create(user *domain.User) error
	GetByEmail(email string) (domain.User, error)
	GetByID(id uint) (domain.User, error)
	Update(user *domain.User) error

	GetSessionByRefreshToken(refreshToken string) (domain.SessionDTO, error)
	SetSession(session *domain.Session) error
	GetSessionByUserID(userID uint) (domain.Session, error)
	DeleteSessionByUserID(userId string) error

	GetAllUsers() ([]domain.User, error)
	DeleteUser(userID string) error
	GetAllSessions() ([]domain.Session, error)
	DeleteSession(sessionID string) error
}

type Room interface {
	Create(room *domain.Room) error
	GetAllRooms() ([]domain.Room, error)
	GetAllMessagesRoom(roomID string) ([]domain.OutputMessageDTO, error)
	CreateMessage(message *domain.Message) error

	UpdateRoom(room *domain.Room) error
	DeleteRoom(roomID string) error
}
