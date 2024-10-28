package main

import (
	"cm/services/gen/authpb"
	"cm/services/sso/config"
	"cm/services/sso/internal/cache"
	"cm/services/sso/internal/logger"
	"cm/services/sso/internal/repository"
	"cm/services/sso/internal/service"
	transport "cm/services/sso/internal/transport/grpc"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	log := logger.New()

	postgresCfg, err := config.LoadPostgresConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed loading cfg")
	}

	redisCfg, err := config.LoadRedisConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed loading cfg")
	}

	serverCfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed loading cfg")
	}

	db, err := config.NewPostgresDb(postgresCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed opening db")
	}

	cacheDB, err := config.NewRedisDb(redisCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed opening cache db")
	}

	repo := repository.New(db)

	cache := cache.New(cacheDB)

	svc := service.New(repo, cache)

	authServer := transport.NewGRPCServer(svc, log)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", serverCfg.PORT))
	if err != nil {
		log.Fatal().Err(err).Msg("failed starting listener")
	}

	srv := grpc.NewServer()

	authpb.RegisterAuthenticationServer(srv, authServer)

	exit := make(chan struct{})
	go func() {
		log.Info().Msgf("Server is now listening on port: %d", serverCfg.PORT)
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
