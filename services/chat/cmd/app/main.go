package main

import (
	"cm/libs/config"
	"cm/libs/log"
	"cm/services/chat/internal/handlers"
	"cm/services/chat/internal/router"
	"cm/services/chat/internal/service"
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

	svc, err := service.New(cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	defer svc.GratefulStop()

	hm := handlers.NewManager(logger, svc)

	r := router.NewRouter(hm)

	quit := make(chan struct{})
	go func() {
		logger.Info().Msgf("Server is now listening on port: %s", cfg.Chat.Port)
		if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Chat.Port), r); err != nil {
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
