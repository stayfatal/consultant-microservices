package main

import (
	"cm/gen/chatpb"
	"cm/libs/config"
	"cm/libs/log"
	"cm/services/matchmaking/internal/repository"
	"cm/services/matchmaking/internal/service"
	transport "cm/services/matchmaking/internal/transport/grpc"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	logger := log.New()

	cfg, err := config.LoadConfigs()
	if err != nil {
		logger.LogFatal(err)
	}

	db, err := config.NewPostgresDB()
	if err != nil {
		logger.LogFatal(err)
	}

	repo := repository.New(db)

	svc, err := service.New(repo, logger)
	if err != nil {
		logger.LogFatal(err)
	}

	chatServer := transport.NewGrpcServer(svc, logger)

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Matchmaking.Port))
	if err != nil {
		logger.LogFatal(err)
	}

	srv := grpc.NewServer()
	defer srv.GracefulStop()

	chatpb.RegisterChatServer(srv, chatServer)

	exit := make(chan struct{})
	go func() {
		logger.Info().Msgf("Server is now listening on port: %s", cfg.Matchmaking.Port)
		if err := srv.Serve(l); err != nil {
			logger.Error().Err(err).Msg("")
			exit <- struct{}{}
		}
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		sig := <-c
		logger.Info().Msg(sig.String())
		exit <- struct{}{}
	}()

	<-exit
}
