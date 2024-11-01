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
	"fmt"
	"net"
	"testing"

	"github.com/google/uuid"
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

	l, err := net.Listen("tcp", ":5001")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	authpb.RegisterAuthenticationServer(srv, authSrv)
	go srv.Serve(l)

	conn, err := grpc.NewClient("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	authClient := authpb.NewAuthenticationClient(conn)

	req := &authpb.RegisterRequest{
		Name:         "test",
		Email:        fmt.Sprintf("test%s@testmail.com", uuid.New().String()),
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

	assert.Equal(t, req.Email, claims.Email)
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

	l, err := net.Listen("tcp", ":5000")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	authpb.RegisterAuthenticationServer(srv, authSrv)
	go srv.Serve(l)

	conn, err := grpc.NewClient("localhost:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	authClient := authpb.NewAuthenticationClient(conn)

	testEmail := fmt.Sprintf("test%s@testmail.com", uuid.New().String())
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
	assert.Equal(t, testEmail, gotClaims.Email)

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
	assert.Equal(t, testEmail, expClaims.Email)
}
