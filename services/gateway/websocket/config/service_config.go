package config

import (
	"cm/internal/utils"

	"github.com/spf13/viper"
)

type ServiceConfig struct {
	CHAT_HOST string
	CHAT_PORT int
}

func LoadServiceConfig() (*ServiceConfig, error) {
	path, err := utils.GetPath("services/gateway/websocket/config/service_config.env")
	if err != nil {
		return nil, err
	}
	viper.SetConfigFile(path)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &ServiceConfig{
		CHAT_HOST: viper.GetString("CHAT_HOST"),
		CHAT_PORT: viper.GetInt("CHAT_PORT"),
	}, nil
}
