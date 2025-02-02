package v1

import (
	"github.com/Nikita-Mihailuk/gochat/server/internal/service"
	"github.com/Nikita-Mihailuk/gochat/server/pkg/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HandlerV1HTTP struct {
	services service.Services
	logger   *zap.Logger
}

func NewHandlerV1HTTP(services service.Services) *HandlerV1HTTP {
	return &HandlerV1HTTP{
		services: services,
		logger:   logging.GetLogger(),
	}
}

func (h *HandlerV1HTTP) InitHandlerV1HTTP(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	h.RegisterUserRouts(v1)
	h.RegisterRoomsRouts(v1)
	h.RegisterAdminRouts(v1)
}
