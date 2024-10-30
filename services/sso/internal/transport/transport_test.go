package transport

import (
	"cm/gen/authpb"
	"cm/internal/publicauth"
	"cm/services/sso/internal/cache"
	"cm/services/sso/internal/logger"
	"cm/services/sso/internal/repository"
	"cm/services/sso/internal/service"
	"cm/services/sso/internal/testhelpers"
	transport "cm/services/sso/internal/transport/grpc"
	"context"
	"net"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestRegister(t *testing.T) {
	postgresDB, err := testhelpers.PreparePostgres(t)
	if err != nil {
		t.Fatal(err)
	}

	redisDB, err := testhelpers.PrepareRedis(t)
	if err != nil {
		t.Fatal(err)
	}

	cache := cache.New(redisDB)

	repo := repository.New(postgresDB)

	svc := service.New(repo, cache)

	log := logger.New()

	authSrv := transport.NewGRPCServer(svc, log)

	srv := grpc.NewServer()

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	authpb.RegisterAuthenticationServer(srv, authSrv)
	go srv.Serve(l)

	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	authClient := authpb.NewAuthenticationClient(conn)

	// expectedId := 1
	req := &authpb.RegisterRequest{
		Name:         "test",
		Email:        "test@testmail.com",
		Password:     "123",
		IsConsultant: false,
	}

	resp, err := authClient.Register(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Error != "" {
		t.Fatal(resp.Error)
	}

	claims, err := publicauth.ValidateToken(resp.Token)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, claims)

	// assert.Equal(t, claims.Id, expectedId)
}

func TestLogin(t *testing.T) {
	postgresDB, err := testhelpers.PreparePostgres(t)
	if err != nil {
		t.Fatal(err)
	}

	redisDB, err := testhelpers.PrepareRedis(t)
	if err != nil {
		t.Fatal(err)
	}

	cache := cache.New(redisDB)

	repo := repository.New(postgresDB)

	svc := service.New(repo, cache)

	log := logger.New()

	authSrv := transport.NewGRPCServer(svc, log)

	srv := grpc.NewServer()

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	authpb.RegisterAuthenticationServer(srv, authSrv)
	go srv.Serve(l)

	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	authClient := authpb.NewAuthenticationClient(conn)

	testEmail := "test@testmail.com"
	testPass := "123"
	expResp, err := authClient.Register(context.Background(), &authpb.RegisterRequest{
		Name:         "test",
		Email:        testEmail,
		Password:     testPass,
		IsConsultant: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	if expResp.Error != "" {
		t.Fatal(expResp.Error)
	}

	gotClaims, err := publicauth.ValidateToken(expResp.Token)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, gotClaims)

	gotResp, err := authClient.Login(context.Background(), &authpb.LoginRequest{
		Email:    testEmail,
		Password: testPass,
	})
	if err != nil {
		t.Fatal(err)
	}

	if gotResp.Error != "" {
		t.Fatal(expResp.Error)
	}

	expClaims, err := publicauth.ValidateToken(gotResp.Token)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, expClaims)
	// assert.Equal(t, expClaims.Id, gotClaims.Id)
}
