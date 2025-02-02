package service

import (
	"github.com/Nikita-Mihailuk/gochat/server/internal/domain"
)

type User interface {
	RegisterUserService(input domain.InputUserDTO) error
	LoginUserService(input domain.InputUserDTO) (domain.Tokens, error)
	GetProfileService(userID uint) (domain.User, error)
	UpdateProfileService(update domain.UpdateProfileDTO) (domain.User, error)

	RefreshTokens(refreshToken string) (domain.Tokens, error)
	DeleteSessionServiceByUserID(userID string) error
}

type Rooms interface {
	GetRoomsService() ([]domain.Room, error)
	CreateRoomService(input domain.InputRoomDTO) error
	GetRoomMessageService(roomId string) ([]domain.OutputMessageDTO, error)
	CreateMessageService(input domain.InputMessageDTO) error
}

type Admin interface {
	GetUsersService() ([]domain.User, error)
	DeleteUserService(userId string) error
	UpdateUserService(update domain.UpdateUserDTO) error

	UpdateRoomService(input domain.Room) error
	DeleteRoomService(roomID string) error

	GetSessionsService() ([]domain.Session, error)
	DeleteSessionService(sessionID string) error
}

type Services struct {
	User  User
	Rooms Rooms
	Admin Admin
}
