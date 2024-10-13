package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBSslMode  string
	Port       int
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("config.env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "reading config")
	}

	cfg := &Config{
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		DBSslMode:  viper.GetString("DB_SSL_MODE"),
		Port:       viper.GetInt("PORT"),
	}

	return cfg, nil
}
