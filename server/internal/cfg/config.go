package cfg

import (
	"github.com/Nikita-Mihailuk/gochat/server/pkg/logging"
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	Listen struct {
		BindIP string `yaml:"bindIp"`
		Port   string `yaml:"port"`
	} `yaml:"listen"`
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"userName"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbName"`
	} `yaml:"database"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Read application configuration")
		instance = &Config{}

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./internal/cfg")

		if err := viper.ReadInConfig(); err != nil {
			logger.Fatal(err.Error())
		}

		if err := viper.Unmarshal(instance); err != nil {
			logger.Fatal("unable to decode cfg into struct: " + err.Error())
		}
	})
	return instance
}
