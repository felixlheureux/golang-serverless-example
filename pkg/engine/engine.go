package engine

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/manta-coder/golang-serverless-example/pkg/db"
	"github.com/manta-coder/golang-serverless-example/pkg/helpers"
	"github.com/manta-coder/golang-serverless-example/pkg/server"
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
	FrontEndDomain                 string `env:"FRONT_END_DOMAIN"`
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

	e := server.NewEcho(logger, config.FrontEndDomain)

	return &Server{e, logger, sql}
}
