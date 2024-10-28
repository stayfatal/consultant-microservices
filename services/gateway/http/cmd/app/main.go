package main

import (
	"cm/services/gateway/http/internal/logger"
	"cm/services/gateway/http/internal/service"
	transport "cm/services/gateway/http/internal/transport/http"
	"net/http"
)

func main() {
	log := logger.New()

	svc, err := service.New()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	srv := transport.NewGatewayServer(svc)
	http.ListenAndServe(":3000", srv)
}
