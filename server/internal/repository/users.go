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
