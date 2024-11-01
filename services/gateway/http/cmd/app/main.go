package main

import (
	"cm/internal/log"
	"cm/services/gateway/http/config"
	"cm/services/gateway/http/internal/service"
	transport "cm/services/gateway/http/internal/transport/http"
	"fmt"
	"net"
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

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", serverCfg.PORT))
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	defer l.Close()

	quit := make(chan struct{})
	go func() {
		logger.Info().Msgf("Server is now listening on port: %d", serverCfg.PORT)
		if err := http.Serve(l, srv); err != nil {
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
