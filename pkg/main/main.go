package main

import (
	"github.com/childrenofukiyo/odin/pkg/db"
	"github.com/childrenofukiyo/odin/pkg/helpers"
	"github.com/childrenofukiyo/odin/pkg/server"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Config struct {
	DBHost                         string `env:"DB_HOST"`
	DBPort                         string `env:"DB_PORT"`
	DBName                         string `env:"DB_NAME"`
	DBUser                         string `env:"DB_USER"`
	DBPass                         string `env:"DB_PASS"`
	LogsDebug                      bool   `env:"LOGS_DEBUG"`
	AuthTokenExpiryDurationSeconds int    `env:"AUTH_TOKEN_EXPIRY_DURATION_SECONDS"`
	AuthSecret                     string `env:"AUTH_SECRET"`
	SanctuaryDomain                string `env:"SANCTUARY_DOMAIN"`
	DopplerEnvironment             string `env:"DOPPLER_ENVIRONMENT"`
}

type Server struct {
	Echo   *echo.Echo
	Logger *zap.SugaredLogger
	DB     *sqlx.DB
}

// MustServer creates a server or panics if there's an error
func MustServer(config Config) *Server {
	// initialize loggers
	logger := helpers.NewLogger(config.LogsDebug)

	sql, err := db.Postgres(config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)
	if err != nil {
		logger.Fatalw("unable to connect to database", "err", err)
	}

	e := server.NewEcho(logger, config.SanctuaryDomain)

	return &Server{e, logger, sql}
}
