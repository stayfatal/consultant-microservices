package transport

import (
	"cm/services/sso/internal/auth"
	"cm/services/sso/internal/logger"
	"cm/services/sso/internal/repository"
	"cm/services/sso/internal/service"
	transport "cm/services/sso/internal/transport/grpc"
	"cm/services/sso/internal/transport/pb"
	"cm/services/sso/internal/utils"
	"context"
	"net"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestRegister(t *testing.T) {
	db, err := utils.PrepareTestingDB()
	if err != nil {
		t.Fatal(err)
	}
	defer utils.ClearTestingDB(db)

	repo := repository.New(db)

	svc := service.New(repo)

	log := logger.New()

	authSrv := transport.NewGRPCServer(svc, log)

	srv := grpc.NewServer()

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	pb.RegisterAuthenticationServer(srv, authSrv)
	go srv.Serve(l)

	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	authClient := pb.NewAuthenticationClient(conn)

	expectedId := 1
	req := &pb.RegisterRequest{
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

	claims, err := auth.ValidateToken(resp.Token)
	if err != nil {
		t.Fatal(err)
	}

	if claims.Id != expectedId {
		t.Errorf("expected id in token %d got %d", expectedId, claims.Id)
	}
}

func TestLogin(t *testing.T) {
	db, err := utils.PrepareTestingDB()
	if err != nil {
		t.Fatal(err)
	}
	defer utils.ClearTestingDB(db)

	repo := repository.New(db)

	svc := service.New(repo)

	log := logger.New()

	authSrv := transport.NewGRPCServer(svc, log)

	srv := grpc.NewServer()

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	pb.RegisterAuthenticationServer(srv, authSrv)
	go srv.Serve(l)

	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	authClient := pb.NewAuthenticationClient(conn)

	testEmail := "test@testmail.com"
	testPass := "123"
	expResp, err := authClient.Register(context.Background(), &pb.RegisterRequest{
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

	gotClaims, err := auth.ValidateToken(expResp.Token)
	if err != nil {
		t.Fatal(err)
	}

	gotResp, err := authClient.Login(context.Background(), &pb.LoginRequest{
		Email:    testEmail,
		Password: testPass,
	})
	if err != nil {
		t.Fatal(err)
	}

	if gotResp.Error != "" {
		t.Fatal(expResp.Error)
	}

	expClaims, err := auth.ValidateToken(gotResp.Token)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expResp.Token, gotResp.Token)
	assert.Equal(t, expClaims.Id, gotClaims.Id)
}
