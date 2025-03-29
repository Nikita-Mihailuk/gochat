package service

import (
	"fmt"
	"github.com/Nikita-Mihailuk/gochat/server/internal/domain"
	"github.com/Nikita-Mihailuk/gochat/server/internal/repository"
	"github.com/Nikita-Mihailuk/gochat/server/internal/service/token_manager"
	"github.com/Nikita-Mihailuk/gochat/server/pkg/auth"
	"github.com/Nikita-Mihailuk/gochat/server/pkg/logging"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"time"
)

type usersService struct {
	userRepo        repository.User
	sessionRepo     repository.Session
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	logger          *zap.Logger
}

func NewUsersService(userRepo repository.User, sessionRepo repository.Session, refreshTokenTTL time.Duration, accessTokenTTL time.Duration) User {
	return &usersService{
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
		tokenManager:    token_manager.GetTokenManager(),
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		logger:          logging.GetLogger(),
	}
}

func (s *usersService) RegisterUserService(input domain.InputUserDTO) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := domain.User{Email: input.Email, Name: input.Name, PasswordHash: string(hashedPassword)}

	err := s.userRepo.Create(&user)
	if err != nil {
		return fmt.Errorf("Пользователь с такой почтой уже существует")
	}
	return nil
}

func (s *usersService) LoginUserService(input domain.InputUserDTO) (domain.Tokens, error) {
	user, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		return domain.Tokens{}, fmt.Errorf("Неверный логин или пароль")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return domain.Tokens{}, fmt.Errorf("Неверный логин или пароль")
	}

	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return domain.Tokens{}, fmt.Errorf("Ошибка при создании нового refresh токена")
	}
	session, err := s.sessionRepo.GetByUserID(user.ID)
	if err != nil {
		session.UserID = user.ID
	}
	session.RefreshToken = refreshToken
	session.ExpiresAt = time.Now().Add(s.accessTokenTTL)

	err = s.sessionRepo.Set(&session)
	if err != nil {
		return domain.Tokens{}, fmt.Errorf("Ошибка при создании сессии")
	}

	accessToken, err := s.tokenManager.NewJWT(strconv.Itoa(int(user.ID)), user.Role, s.accessTokenTTL)
	if err != nil {
		return domain.Tokens{}, fmt.Errorf("Ошибка при создании нового access токен")
	}
	return domain.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *usersService) GetProfileService(id uint) (domain.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *usersService) UpdateProfileService(update domain.UpdateProfileDTO) (domain.User, error) {
	user, err := s.userRepo.GetByID(update.UserId)
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
		if err = SaveFile(update.FileHeader, filePath); err != nil {
			return domain.User{}, fmt.Errorf("Ошибка сохранения фото")
		}
		user.PhotoURL = filePath
	}
	user.UpdatedAt = time.Now()

	err = s.userRepo.Update(&user)
	if err != nil {
		return domain.User{}, fmt.Errorf("Ошибка обновления профиля")
	}
	return user, nil
}

func (s *usersService) RefreshTokens(refreshToken string) (domain.Tokens, error) {
	session, err := s.sessionRepo.GetByRefreshToken(refreshToken)
	if err != nil {
		return domain.Tokens{}, fmt.Errorf("Ошибка при получении сессии пользователя")
	}

	newRefreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return domain.Tokens{}, fmt.Errorf("Ошибка при обновлении refresh токена")
	}

	newAccessToken, err := s.tokenManager.NewJWT(strconv.Itoa(int(session.UserID)), session.UserRole, s.accessTokenTTL)
	if err != nil {
		return domain.Tokens{}, fmt.Errorf("Ошибка при обновлении access токена")
	}

	err = s.sessionRepo.Set(&domain.Session{
		ID:           session.ID,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
		UserID:       session.UserID,
	})
	if err != nil {
		return domain.Tokens{}, fmt.Errorf("Ошибка при обновлении сессии")
	}

	return domain.Tokens{AccessToken: newAccessToken, RefreshToken: newRefreshToken}, nil
}
func (s *usersService) DeleteSessionServiceByUserID(userID string) error {
	err := s.sessionRepo.DeleteByUserID(userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *usersService) GetTokenManager() auth.TokenManager {
	return s.tokenManager
}

func SaveFile(fileHeader *multipart.FileHeader, path string) error {
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
