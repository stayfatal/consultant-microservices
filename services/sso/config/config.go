package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_NAME     string
	POSTGRES_SSL_MODE string
	POSTGRES_PORT     int
	REDIS_ADDR        string
	REDIS_PASSWORD    string
}

func LoadConfig() (*Config, error) {
	path, err := getFilePath("config.env")
	if err != nil {
		return nil, err
	}
	viper.SetConfigFile(path)
	err = viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "reading config")
	}

	cfg := &Config{
		POSTGRES_USER:     viper.GetString("POSTGRES_USER"),
		POSTGRES_PASSWORD: viper.GetString("POSTGRES_PASSWORD"),
		POSTGRES_NAME:     viper.GetString("POSTGRES_NAME"),
		POSTGRES_SSL_MODE: viper.GetString("POSTGRES_SSL_MODE"),
		POSTGRES_PORT:     viper.GetInt("POSTGRES_PORT"),
		REDIS_ADDR:        viper.GetString("REDIS_ADDR"),
		REDIS_PASSWORD:    viper.GetString("REDIS_PASSWORD"),
	}

	return cfg, nil
}

func getFilePath(filename string) (string, error) {
	_, currentFile, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("unable to get caller information")
	}

	dir := filepath.Dir(currentFile)

	fullPath := filepath.Join(dir, filename)

	return fullPath, nil
}
