package repository

import "github.com/Nikita-Mihailuk/gochat/server/internal/domain"

type User interface {
	Create(user *domain.User) error
	GetByEmail(email string) (domain.User, error)
	GetByID(id uint) (domain.User, error)
	Update(user *domain.User) error
	GetAllUsers() ([]domain.User, error)
	Delete(userID string) error
}

type Room interface {
	Create(room *domain.Room) error
	GetAll() ([]domain.Room, error)
	GetAllMessages(roomID string) ([]domain.OutputMessageDTO, error)
	CreateMessage(message *domain.Message) error
	Update(room *domain.Room) error
	Delete(roomID string) error
}

type Session interface {
	GetByRefreshToken(refreshToken string) (domain.SessionDTO, error)
	Set(session *domain.Session) error
	GetByUserID(userID uint) (domain.Session, error)
	DeleteByUserID(userId string) error
	GetAll() ([]domain.Session, error)
	DeleteByID(sessionID string) error
}
