package dbClients

import (
	"fmt"
	"github.com/Nikita-Mihailuk/gochat/server/internal/cfg"
	"github.com/Nikita-Mihailuk/gochat/server/internal/domain"
	"github.com/Nikita-Mihailuk/gochat/server/pkg/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func init() {
	var err error

	getConfig := cfg.GetConfig()
	logger := logging.GetLogger()

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		getConfig.Database.Host,
		getConfig.Database.Username,
		getConfig.Database.Password,
		getConfig.Database.DBName,
		getConfig.Database.Port)

	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		logger.Fatal("Ошибка подключения к базе данных: " + err.Error())
	}

	err = db.AutoMigrate(&domain.User{}, &domain.Room{}, &domain.Message{}, &domain.Session{})
	if err != nil {
		logger.Fatal(err.Error())
	}
}
