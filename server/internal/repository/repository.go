package repository

import "github.com/Nikita-Mihailuk/gochat/server/internal/domain"

type User interface {
	Create(user *domain.User) error
	GetByEmail(email string) (domain.User, error)
	GetByID(id uint) (domain.User, error)
	Update(user *domain.User) error
}

type Room interface {
	Create(room *domain.Room) error
	GetAllRooms() ([]domain.Room, error)
	GetAllMessagesRoom(roomID string) ([]domain.OutputMessageDTO, error)
	CreateMessage(message *domain.Message) error
}
