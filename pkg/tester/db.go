package tester

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/manta-coder/golang-serverless-example/pkg/db"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
	"io"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

const (
	postgresImage           = "postgres:12.7"
	postgresUser            = "ukiyo"
	postgresPassword        = "pass123"
	postgresDefaultDatabase = "odin"
)

type Database string

var dbOnce sync.Once
var dbInstance *sqlx.DB

func DB() *sqlx.DB {
	dbOnce.Do(initDB)

	return dbInstance
}

// initDB a postgres database via testcontainers
func initDB() {
	ctx := context.Background()

	var host string

	host = "localhost"

	url := func(port nat.Port) string {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPassword, host, port.Port(), postgresDefaultDatabase)
	}

	req := testcontainers.ContainerRequest{
		Image:        postgresImage,
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForSQL("5432/tcp", "postgres", url).Timeout(1 * time.Minute),
		Env: map[string]string{
			"POSTGRES_USER":     postgresUser,
			"POSTGRES_PASSWORD": postgresPassword,
			"POSTGRES_DB":       postgresDefaultDatabase,
		},
		ReaperImage: testcontainersReaper,
	}

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		logger.Panic("can't create postgres postgres", zap.Error(err))
	}

	port, err := postgres.MappedPort(ctx, "5432/tcp")

	if err != nil {
		logger.Panic("can't get mapped port postgres", zap.Error(err))
	}

	logger.Info("started postgres container",
		zap.String("port", port.Port()),
		zap.String("user", postgresUser),
		zap.String("pass", postgresPassword))

	dbInstance = mustMigrate(host, port.Port())
}

// mustMigrate runs a migrations and returns a connection
func mustMigrate(host string, port string) *sqlx.DB {
	// find the current directory so we can infer the migrations folder
	_, filename, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(filename)
	if err != nil {
		logger.Panic(zap.Error(err))
	}

	cmdName := "bash"
	migrationScriptPath := filepath.Join(dir, "../../../scripts/migrate.sh")

	// on windows, we must execute the script with git bash
	if runtime.GOOS == "windows" {
		cmdName = "C:\\Program Files\\Git\\bin\\bash.exe"
	}

	cmd := exec.Command(cmdName, migrationScriptPath, "--db-host", host, "--db-port", port)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		logger.Panic(zap.Error(err))
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()

	if err != nil {
		logger.Panic(zap.Error(err))
	}
	defer stderr.Close()

	err = cmd.Start()

	if err != nil {
		logger.Panic(zap.Error(err))
	}

	scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for scanner.Scan() {
		logger.Info(scanner.Text())
	}

	err = cmd.Wait()

	if err != nil {
		logger.Panic(zap.Error(err))
	}

	return mustConnectDB(host, port, postgresDefaultDatabase)
}

// mustConnectDB creates a db connection or panics
func mustConnectDB(host, port, database string) *sqlx.DB {
	d, err := db.Postgres(host, port, postgresUser, postgresPassword, database)

	if err != nil {
		logger.Panic("unable connect to db", zap.Error(err))
	}

	return d
}
