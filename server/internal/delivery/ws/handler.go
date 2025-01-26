package ws

import (
	"github.com/Nikita-Mihailuk/gochat/server/internal/delivery/ws/v1"
	"github.com/Nikita-Mihailuk/gochat/server/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services service.Services
}

func NewHandlerWS(services service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitHandlerWS(router *gin.Engine) {
	handlerV1 := v1.NewHandlerV1WS(h.services)
	api := router.Group("/api")
	handlerV1.InitHandlerV1WS(api)
}
