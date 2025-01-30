package token_manager

import (
	"github.com/Nikita-Mihailuk/gochat/server/internal/cfg"
	"github.com/Nikita-Mihailuk/gochat/server/pkg/auth"
	"github.com/Nikita-Mihailuk/gochat/server/pkg/logging"
	"sync"
)

var tokenManager auth.TokenManager
var err error
var once sync.Once

func GetTokenManager() auth.TokenManager {
	once.Do(func() {
		logger := logging.GetLogger()
		config := cfg.GetConfig()
		tokenManager, err = auth.NewManager(config.Auth.SecretKey)
		if err != nil {
			logger.Fatal(err.Error())
		}
	})
	return tokenManager
}
