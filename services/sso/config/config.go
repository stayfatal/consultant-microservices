package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var Cfg *Config

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBSslMode  string
	JWTSecret  []byte
	Port       int
}

func LoadConfig() error {
	viper.SetConfigFile("config.env")
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "reading config")
	}

	Cfg = &Config{
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		DBSslMode:  viper.GetString("DB_SSL_MODE"),
		JWTSecret:  []byte(viper.GetString("JWT_SECRET")),
		Port:       viper.GetInt("PORT"),
	}

	return nil
}
