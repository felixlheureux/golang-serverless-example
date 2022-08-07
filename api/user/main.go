package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/caarlos0/env/v6"
	"github.com/labstack/echo/v4"
	"github.com/manta-coder/golang-serverless-example/pkg/controller"
	"github.com/manta-coder/golang-serverless-example/pkg/engine"
	"github.com/manta-coder/golang-serverless-example/pkg/service"
	"github.com/manta-coder/golang-serverless-example/pkg/store"
)

var echoLambda *echoadapter.EchoLambdaV2

func init() {
	var config engine.Config
	if err := env.Parse(&config); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	server := engine.MustServer(config)

	userStore := store.NewUserStore(server.Logger, server.DB)
	userService := service.NewUserService(server.Logger, userStore)

	middlewares := []echo.MiddlewareFunc{
		controller.NewAuthenticator(config.AuthSecret),
		controller.NewAuthMiddleware(),
	}

	group := server.Echo.Group("/users", middlewares...)
	controller.NewUserController(group, server.Logger, userService)

	echoLambda = echoadapter.NewV2(server.Echo)
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return echoLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
