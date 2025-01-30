package service

import (
	"github.com/Nikita-Mihailuk/gochat/server/internal/domain"
)

type User interface {
	RegisterUserService(input domain.InputUserDTO) error
	LoginUserService(input domain.InputUserDTO) (domain.Tokens, error)
	GetProfileService(id uint) (domain.User, error)
	UpdateProfileService(update domain.UpdateUserDTO) (domain.User, error)

	RefreshTokens(refreshToken string) (domain.Tokens, error)
	DeleteSessionService(userID string) error
}

type Rooms interface {
	GetRoomsService() ([]domain.Room, error)
	CreateRoomService(input domain.InputRoomDTO) error
	GetRoomMessageService(roomId string) ([]domain.OutputMessageDTO, error)
	CreateMessageService(input domain.InputMessageDTO) error
}

type Services struct {
	User  User
	Rooms Rooms
}
