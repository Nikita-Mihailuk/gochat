package repository

import (
	"github.com/Nikita-Mihailuk/gochat/server/internal/domain"
	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) Session {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) GetByRefreshToken(refreshToken string) (domain.SessionDTO, error) {
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

func (r *sessionRepository) Set(session *domain.Session) error {
	err := r.db.Save(session).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *sessionRepository) GetByUserID(userID uint) (domain.Session, error) {
	var session domain.Session
	err := r.db.Where("user_id = ?", userID).First(&session).Error
	if err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

func (r *sessionRepository) DeleteByUserID(userID string) error {
	err := r.db.Where("user_id = ?", userID).Delete(&domain.Session{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *sessionRepository) GetAll() ([]domain.Session, error) {
	var sessions []domain.Session
	err := r.db.Order("id").Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *sessionRepository) DeleteByID(sessionID string) error {
	err := r.db.Delete(&domain.Session{}, sessionID).Error
	if err != nil {
		return err
	}
	return nil
}
