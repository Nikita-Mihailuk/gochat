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

func (r *usersRepository) GetSessionByRefreshToken(refreshToken string) (domain.SessionDTO, error) {
	var session domain.SessionDTO
	err := r.db.Table("sessions").
		Select("sessions.id, sessions.refresh_token, sessions.expires_at, sessions.user_id, users.role as user_role").
		Joins("JOIN users ON users.id = sessions.user_id").
		Where("sessions.refresh_token = ?", refreshToken).
		Scan(&session).Error
	if err != nil {
		return domain.SessionDTO{}, err
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

func (r *usersRepository) GetSessionByUserID(userID uint) (domain.Session, error) {
	var session domain.Session
	err := r.db.Where("user_id = ?", userID).First(&session).Error
	if err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

func (r *usersRepository) DeleteSessionByUserID(userID string) error {
	err := r.db.Where("user_id = ?", userID).Delete(&domain.Session{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *usersRepository) GetAllUsers() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Not("role = ?", "admin").Order("id").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *usersRepository) DeleteUser(userID string) error {
	err := r.db.Delete(&domain.User{}, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *usersRepository) GetAllSessions() ([]domain.Session, error) {
	var sessions []domain.Session
	err := r.db.Order("id").Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *usersRepository) DeleteSession(sessionID string) error {
	err := r.db.Delete(&domain.Session{}, sessionID).Error
	if err != nil {
		return err
	}
	return nil
}
