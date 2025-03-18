package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	test := map[string]string{
		"AUTH_HOST":        "auth",
		"AUTH_PORT":        "8030",
		"GATEWAY_HOST":     "gateway",
		"GATEWAY_PORT":     "8020",
		"MATCHMAKING_HOST": "matchmaking",
		"MATCHMAKING_PORT": "8010",
		"CHAT_HOST":        "chat",
		"CHAT_PORT":        "8000",
	}

	for key, val := range test {
		os.Setenv(key, val)
	}

	cfg, err := LoadConfigs()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "auth", cfg.Auth.Host)
	assert.Equal(t, "8030", cfg.Auth.Port)

	assert.Equal(t, "gateway", cfg.Gateway.Host)
	assert.Equal(t, "8020", cfg.Gateway.Port)

	assert.Equal(t, "matchmaking", cfg.Matchmaking.Host)
	assert.Equal(t, "8010", cfg.Matchmaking.Port)

	assert.Equal(t, "chat", cfg.Chat.Host)
	assert.Equal(t, "8000", cfg.Chat.Port)
}
