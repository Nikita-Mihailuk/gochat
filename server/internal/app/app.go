package app

import (
	"fmt"
	"github.com/Nikita-Mihailuk/gochat/server/internal/cfg"
	httpp "github.com/Nikita-Mihailuk/gochat/server/internal/delivery/http"
	"github.com/Nikita-Mihailuk/gochat/server/internal/delivery/ws"
	"github.com/Nikita-Mihailuk/gochat/server/internal/repository"
	"github.com/Nikita-Mihailuk/gochat/server/internal/service"
	"github.com/Nikita-Mihailuk/gochat/server/middleware"
	"github.com/Nikita-Mihailuk/gochat/server/pkg/dbClients"
	"github.com/Nikita-Mihailuk/gochat/server/pkg/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

func Run() {
	logger := logging.GetLogger()
	db := dbClients.GetDB()

	logger.Info("Создание роутера")
	router := gin.Default()
	router.Use(middleware.CorsMiddleware())
	router.Static("/uploads", "./uploads")

	userRepo := repository.NewUsersRepository(db)
	roomRepo := repository.NewRoomsRepository(db)

	services := service.Services{
		User:  service.NewUsersService(userRepo),
		Rooms: service.NewRoomsService(roomRepo),
	}

	httpHandler := httpp.NewHandlerHTTP(services)
	wsHandler := ws.NewHandlerWS(services)

	logger.Info("Регистрация эндпоинтов")
	httpHandler.InitHandlerHTTP(router)
	wsHandler.InitHandlerWS(router)

	start(router)
}

func start(router *gin.Engine) {

	getConfig := cfg.GetConfig()
	logger := logging.GetLogger()

	logger.Info("Прослушивание сервера",
		zap.String("bind_ip", getConfig.Listen.BindIP),
		zap.String("port", getConfig.Listen.Port))

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", getConfig.Listen.BindIP, getConfig.Listen.Port))

	if err != nil {
		logger.Fatal(err.Error())
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Minute,
		ReadTimeout:  15 * time.Minute,
	}

	logger.Fatal(server.Serve(listener).Error())
}
