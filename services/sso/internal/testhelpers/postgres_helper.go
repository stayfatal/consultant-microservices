package testhelpers

import (
	"cm/services/sso/config"
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func ConfigurePostgresContainer(ctx context.Context) (testcontainers.Container, *sqlx.DB, error) {
	containerCfg, err := config.LoadPostgresConfig()
	if err != nil {
		return nil, nil, err
	}

	req := testcontainers.ContainerRequest{
		Image:        containerCfg.IMAGE,
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", containerCfg.PORT)},
		Env: map[string]string{
			"POSTGRES_USER":     containerCfg.USER,
			"POSTGRES_PASSWORD": containerCfg.PASSWORD,
			"POSTGRES_DB":       containerCfg.DB_NAME,
		},
		WaitingFor: wait.ForListeningPort(nat.Port(fmt.Sprintf("%d/tcp", containerCfg.PORT))).WithStartupTimeout(60 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	dbConnCfg := containerCfg

	dbConnCfg.HOST, err = container.Host(ctx)
	if err != nil {
		return nil, nil, err
	}

	port, err := container.MappedPort(ctx, nat.Port(fmt.Sprintf("%d", containerCfg.PORT)))
	if err != nil {
		return nil, nil, err
	}

	dbConnCfg.PORT, err = strconv.Atoi(port.Port())
	if err != nil {
		return nil, nil, err
	}

	db, err := config.NewPostgresDb(containerCfg)
	if err != nil {
		return nil, nil, err
	}

	return container, db, nil
}

func CleanupPostgresContainer(t *testing.T, container testcontainers.Container, db *sqlx.DB) {
	defer func() {
		err := testcontainers.TerminateContainer(container)
		if err != nil {
			t.Fatal(err)
		}
	}()
	err := db.Close()
	if err != nil {
		t.Fatal(err)
	}
}
