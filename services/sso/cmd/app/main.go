package main

import (
	"cm/services/sso/config"
	"cm/services/sso/internal/logger"
	"cm/services/sso/internal/repository"
	"cm/services/sso/internal/service"
	transport "cm/services/sso/internal/transport/grpc"
	"cm/services/sso/internal/transport/pb"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	log := logger.New()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed loading cfg")
	}

	db, err := config.NewDb(*cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed opening db")
	}

	repo := repository.New(db)

	svc := service.New(repo)

	authServer := transport.NewGRPCServer(svc, log)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("failed starting listener")
	}

	srv := grpc.NewServer()

	pb.RegisterAuthenticationServer(srv, authServer)

	exit := make(chan struct{})
	go func() {
		log.Info().Msgf("Server is now listening on port: %d", cfg.Port)
		if err := srv.Serve(l); err != nil {
			log.Error().Err(err).Msg("")
			exit <- struct{}{}
		}
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		sig := <-c
		log.Info().Msg(sig.String())
		exit <- struct{}{}
	}()

	<-exit
	db.Close()
	srv.GracefulStop()
}
