package service

import (
	"fmt"
	"github.com/Nikita-Mihailuk/gochat/server/internal/domain"
	"github.com/Nikita-Mihailuk/gochat/server/internal/repository"
	"github.com/Nikita-Mihailuk/gochat/server/pkg/logging"
	"go.uber.org/zap"
	"time"
)

type adminsService struct {
	userRepo    repository.User
	roomRepo    repository.Room
	sessionRepo repository.Session
	logger      *zap.Logger
}

func NewAdminsService(userRepo repository.User, sessionRepo repository.Session, roomRepo repository.Room) Admin {
	return &adminsService{
		userRepo:    userRepo,
		roomRepo:    roomRepo,
		sessionRepo: sessionRepo,
		logger:      logging.GetLogger(),
	}
}

func (s *adminsService) GetUsersService() ([]domain.User, error) {
	users, err := s.userRepo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("Ошибка при получении списка пользователей")
	}
	return users, nil
}

func (s *adminsService) DeleteUserService(userId string) error {
	err := s.userRepo.Delete(userId)
	if err != nil {
		return fmt.Errorf("Ошибка при удалении пользователя")
	}
	err = s.sessionRepo.DeleteByUserID(userId)
	if err != nil {
		return fmt.Errorf("Ошибка при удалении сессии пользователя")
	}
	return nil
}

func (s *adminsService) GetSessionsService() ([]domain.Session, error) {
	sessions, err := s.sessionRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("Ошибка при получении сессий")
	}
	return sessions, nil
}

func (s *adminsService) UpdateRoomService(input domain.Room) error {
	err := s.roomRepo.Update(&input)
	if err != nil {
		return fmt.Errorf("Комната с таким именем уже существует")
	}
	return nil
}

func (s *adminsService) DeleteRoomService(roomID string) error {
	err := s.roomRepo.Delete(roomID)
	if err != nil {
		return fmt.Errorf("Ошибка при удалении комнаты")
	}
	return nil
}

func (s *adminsService) UpdateUserService(update domain.UpdateUserDTO) error {
	user, err := s.userRepo.GetByID(update.UserId)
	if err != nil {
		return fmt.Errorf("Пользователь не найден")
	}

	if update.NewName != "" {
		user.Name = update.NewName
	}

	if update.FileHeader != nil {
		filePath := fmt.Sprintf("uploads/%d_%s", user.ID, update.FileHeader.Filename)
		if err = SaveFile(update.FileHeader, filePath); err != nil {
			return fmt.Errorf("Ошибка сохранения фото")
		}
		user.PhotoURL = filePath
	}
	user.UpdatedAt = time.Now()

	err = s.userRepo.Update(&user)
	if err != nil {
		return fmt.Errorf("Ошибка обновления профиля")
	}
	return nil
}

func (s *adminsService) DeleteSessionService(sessionID string) error {
	err := s.sessionRepo.DeleteByID(sessionID)
	if err != nil {
		return fmt.Errorf("Ошибка при удалении сессии")
	}
	return nil
}
