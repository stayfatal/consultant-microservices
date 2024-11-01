package config

import (
	"cm/internal/utils"

	"github.com/spf13/viper"
)

type ServiceConfig struct {
	SsoHost string
	SsoPort int
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
		SsoHost: viper.GetString("SSO_HOST"),
		SsoPort: viper.GetInt("SSO_PORT"),
	}, nil
}
