package config

import (
	"cm/internal/utils"

	"github.com/spf13/viper"
)

type ServiceConfig struct {
	SSO_HOST string
	SSO_PORT int
}

func LoadServiceConfig() (*ServiceConfig, error) {
	path, err := utils.GetPath("services/gateway/http/config/service_config.env")
	if err != nil {
		return nil, err
	}
	viper.SetConfigFile(path)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &ServiceConfig{
		SSO_HOST: viper.GetString("SSO_HOST"),
		SSO_PORT: viper.GetInt("SSO_PORT"),
	}, nil
}
