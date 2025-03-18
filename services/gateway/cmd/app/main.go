package main

import (
	"cm/libs/config"
	"cm/libs/log"
	"cm/services/gateway/internal/service"
	transport "cm/services/gateway/internal/transport/http"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.New()

	cfg, err := config.LoadConfigs()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	svc, err := service.New(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	defer svc.GratefulStop()

	srv := transport.NewGatewayServer(svc, logger)

	quit := make(chan struct{})
	go func() {
		logger.Info().Msgf("Server is now listening on port: %s", cfg.Gateway.Port)
		if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Gateway.Port), srv); err != nil {
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
