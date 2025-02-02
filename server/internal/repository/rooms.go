package repository

import (
	"github.com/Nikita-Mihailuk/gochat/server/internal/domain"
	"gorm.io/gorm"
)

type roomsRepository struct {
	db *gorm.DB
}

func NewRoomsRepository(db *gorm.DB) Room {
	return &roomsRepository{db: db}
}

func (r *roomsRepository) Create(room *domain.Room) error {
	err := r.db.Create(room).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *roomsRepository) GetAllRooms() ([]domain.Room, error) {
	var rooms []domain.Room
	err := r.db.Order("id").Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomsRepository) GetAllMessagesRoom(roomID string) ([]domain.OutputMessageDTO, error) {
	var messages []domain.OutputMessageDTO
	err := r.db.Table("messages").
		Select("messages.*, users.photo_url AS user_avatar, users.name AS user_name").
		Joins("JOIN users ON users.id = messages.user_id").
		Where("messages.room_id = ?", roomID).
		Order("messages.created_at").
		Find(&messages).Error

	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *roomsRepository) CreateMessage(message *domain.Message) error {
	err := r.db.Create(message).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *roomsRepository) UpdateRoom(room *domain.Room) error {
	err := r.db.Save(room).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *roomsRepository) DeleteRoom(roomID string) error {
	err := r.db.Delete(domain.Room{}, roomID).Error
	if err != nil {
		return err
	}
	return nil
}
