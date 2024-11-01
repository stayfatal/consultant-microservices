package main

import (
	"cm/services/gateway/http/config"
	"cm/services/gateway/http/internal/logger"
	"cm/services/gateway/http/internal/service"
	transport "cm/services/gateway/http/internal/transport/http"
	"net/http"
)

func main() {
	log := logger.New()

	cfg, err := config.LoadServiceConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	svc, err := service.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	srv := transport.NewGatewayServer(svc)
	http.ListenAndServe(":3000", srv)
}
