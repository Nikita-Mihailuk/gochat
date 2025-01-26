package service

import (
	"fmt"
	"github.com/Nikita-Mihailuk/gochat/server/internal/domain"
	"github.com/Nikita-Mihailuk/gochat/server/internal/repository"
	"github.com/Nikita-Mihailuk/gochat/server/pkg/logging"
	"go.uber.org/zap"
)

type roomsService struct {
	repo   repository.Room
	logger *zap.Logger
}

func NewRoomsService(repo repository.Room) Rooms {
	return &roomsService{
		repo:   repo,
		logger: logging.GetLogger(),
	}
}

func (s *roomsService) GetRoomsService() ([]domain.Room, error) {
	rooms, err := s.repo.GetAllRooms()
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *roomsService) CreateRoomService(input domain.InputRoomDTO) error {
	room := domain.Room{Name: input.Name}
	err := s.repo.Create(&room)
	if err != nil {
		return fmt.Errorf("Комната с таким именем уже существует")
	}
	return nil
}

func (s *roomsService) GetRoomMessageService(roomId string) ([]domain.OutputMessageDTO, error) {
	messages, err := s.repo.GetAllMessagesRoom(roomId)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *roomsService) CreateMessageService(input domain.InputMessageDTO) error {
	err := s.repo.CreateMessage(&domain.Message{RoomID: input.RoomID, UserID: input.UserID, Content: input.Content})
	if err != nil {
		return err
	}
	return nil
}
