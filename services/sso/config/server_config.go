package config

import (
	"cm/internal/utils"

	"github.com/pkg/errors"
	viperr "github.com/spf13/viper"
)

type ServerConfig struct {
	PORT int
}

func LoadServerConfig() (*ServerConfig, error) {
	path, err := utils.GetPath("services/sso/config/server_config.env")
	if err != nil {
		return nil, err
	}
	viperr.SetConfigFile(path)
	err = viperr.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "reading config")
	}

	cfg := &ServerConfig{
		PORT: viperr.GetInt("PORT"),
	}

	return cfg, nil
}
