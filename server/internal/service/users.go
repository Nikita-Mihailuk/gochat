package service

import (
	"fmt"
	"github.com/Nikita-Mihailuk/gochat/server/internal/domain"
	"github.com/Nikita-Mihailuk/gochat/server/internal/repository"
	"github.com/Nikita-Mihailuk/gochat/server/pkg/logging"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"io"
	"mime/multipart"
	"os"
	"time"
)

type usersService struct {
	repo   repository.User
	logger *zap.Logger
}

func NewUsersService(repo repository.User) User {
	return &usersService{
		repo:   repo,
		logger: logging.GetLogger(),
	}
}

func (s *usersService) RegisterUserService(input domain.InputUserDTO) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := domain.User{Email: input.Email, Name: input.Name, PasswordHash: string(hashedPassword)}

	err := s.repo.Create(&user)
	if err != nil {
		return err
	}
	return nil
}

func (s *usersService) LoginUserService(input domain.InputUserDTO) (uint, error) {
	user, err := s.repo.GetByEmail(input.Email)
	if err != nil {
		return 0, fmt.Errorf("Неверный логин или пароль")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return 0, fmt.Errorf("Неверный логин или пароль")
	}
	return user.ID, nil
}

func (s *usersService) GetProfileService(id uint) (domain.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *usersService) UpdateProfileService(update domain.UpdateUserDTO) (domain.User, error) {
	user, err := s.repo.GetByID(update.UserId)
	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(update.CurrentPassword))
	if err != nil {
		return domain.User{}, fmt.Errorf("Неверный текущий пароль")
	}

	if update.NewName != "" {
		user.Name = update.NewName
	}

	if update.NewPassword != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(update.NewPassword), bcrypt.DefaultCost)
		user.PasswordHash = string(hashedPassword)
	}

	if update.FileHeader != nil {
		filePath := fmt.Sprintf("uploads/%d_%s", user.ID, update.FileHeader.Filename)
		if err = s.saveFile(update.FileHeader, filePath); err != nil {
			return domain.User{}, fmt.Errorf("Ошибка сохранения фото")
		}
		user.PhotoURL = filePath
	}
	user.UpdatedAt = time.Now()

	err = s.repo.Update(&user)
	if err != nil {
		return domain.User{}, fmt.Errorf("Ошибка обновления профиля")
	}
	return user, nil
}

func (s *usersService) saveFile(fileHeader *multipart.FileHeader, path string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}
