package config

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type ServicesConfig struct {
	Chat        ServiceConfig
	Gateway     ServiceConfig
	Matchmaking ServiceConfig
	Auth        ServiceConfig
}

type ServiceConfig struct {
	Host string
	Port string
}

func LoadConfigs() (*ServicesConfig, error) {
	viper.AutomaticEnv()

	prefixes := map[string]*ServiceConfig{
		"AUTH":        {},
		"GATEWAY":     {},
		"MATCHMAKING": {},
		"CHAT":        {},
	}

	for prefix, cfg := range prefixes {
		hostKey := fmt.Sprintf("%s_HOST", prefix)
		portKey := fmt.Sprintf("%s_PORT", prefix)

		err := viper.BindEnv(hostKey)
		if err != nil {
			return nil, errors.Wrap(err, "binding env")
		}
		err = viper.BindEnv(portKey)
		if err != nil {
			return nil, errors.Wrap(err, "binding env")
		}

		cfg.Host = viper.GetString(hostKey)
		cfg.Port = viper.GetString(portKey)
	}

	return &ServicesConfig{
		Auth:        *prefixes["AUTH"],
		Gateway:     *prefixes["GATEWAY"],
		Matchmaking: *prefixes["MATCHMAKING"],
		Chat:        *prefixes["CHAT"],
	}, nil
}
