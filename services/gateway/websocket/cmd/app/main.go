package main

import (
	"cm/internal/log"
	"cm/services/gateway/websocket/config"
	"cm/services/gateway/websocket/internal/handlers"
	"cm/services/gateway/websocket/internal/router"
	"cm/services/gateway/websocket/internal/service"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.New()
	cfg, err := config.LoadServerConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	svc := service.New()

	hm := handlers.NewManager(logger, svc)

	r := router.NewRouter(hm)

	quit := make(chan struct{})
	go func() {
		logger.Info().Msgf("Server is now listening on port: %d", cfg.PORT)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.PORT), r); err != nil {
			logger.Fatal().Err(err).Msg("")
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
