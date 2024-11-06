package config

import (
	"cm/internal/utils"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	PORT int
}

func LoadServerConfig() (*ServerConfig, error) {
	path, err := utils.GetPath("services/gateway/websocket/config/server_config.env")
	if err != nil {
		return nil, err
	}
	viper.SetConfigFile(path)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &ServerConfig{
		PORT: viper.GetInt("PORT"),
	}, nil
}
