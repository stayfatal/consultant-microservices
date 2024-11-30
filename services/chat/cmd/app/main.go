package main

import (
	"cm/gen/chatpb"
	"cm/internal/log"
	"cm/services/chat/config"
	"cm/services/chat/internal/repository"
	"cm/services/chat/internal/service"
	transport "cm/services/chat/internal/transport/grpc"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	logger := log.New()

	postgresCfg, err := config.LoadPostgresConfig()
	if err != nil {
		logger.LogFatal(err)
	}

	db, err := config.NewPostgresDb(postgresCfg)
	if err != nil {
		logger.LogFatal(err)
	}

	repo := repository.New(db)

	svc, err := service.New(repo, logger)
	if err != nil {
		logger.LogFatal(err)
	}

	chatServer := transport.NewGrpcServer(svc, logger)

	serverCfg, err := config.LoadServerConfig()
	if err != nil {
		logger.LogFatal(err)
	}

	l, err := net.Listen("tcp", fmt.Sprint(":", serverCfg.PORT))
	if err != nil {
		logger.LogFatal(err)
	}

	srv := grpc.NewServer()
	defer srv.GracefulStop()

	chatpb.RegisterChatServer(srv, chatServer)

	exit := make(chan struct{})
	go func() {
		logger.Info().Msgf("Server is now listening on port: %d", serverCfg.PORT)
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
