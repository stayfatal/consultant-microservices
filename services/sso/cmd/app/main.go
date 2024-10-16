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

	srv.Serve(l)
}
