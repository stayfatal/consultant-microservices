package main

import (
	"cm/internal/log"
	"cm/services/gateway/http/config"
	"cm/services/gateway/http/internal/service"
	transport "cm/services/gateway/http/internal/transport/http"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.New()

	serviceCfg, err := config.LoadServiceConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	svc, err := service.New(serviceCfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	serverCfg, err := config.LoadServerConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	srv := transport.NewGatewayServer(svc, logger)

	quit := make(chan struct{})
	go func() {
		logger.Info().Msgf("Server is now listening on port: %d", serverCfg.PORT)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", serverCfg.PORT), srv); err != nil {
			logger.Error().Err(err).Msg("")
			quit <- struct{}{}
		}
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		sig := <-c
		logger.Info().Msg(sig.String())
		quit <- struct{}{}
	}()

	<-quit
}
