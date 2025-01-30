package repository

import (
	"github.com/Nikita-Mihailuk/gochat/server/internal/domain"
	"gorm.io/gorm"
)

type usersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) User {
	return &usersRepository{db: db}
}

func (r *usersRepository) Create(user *domain.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *usersRepository) GetByEmail(email string) (domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *usersRepository) GetByID(id uint) (domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *usersRepository) Update(user *domain.User) error {
	err := r.db.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *usersRepository) GetByRefreshToken(refreshToken string) (domain.Session, error) {
	var session domain.Session
	err := r.db.Where("refresh_token = ?", refreshToken).First(&session).Error
	if err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

func (r *usersRepository) SetSession(session *domain.Session) error {
	err := r.db.Save(session).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *usersRepository) DeleteSession(userID string) error {
	err := r.db.Where("user_id = ?", userID).Delete(&domain.Session{}).Error
	if err != nil {
		return err
	}
	return nil
}
