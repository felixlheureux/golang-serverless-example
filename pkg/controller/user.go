package controller

import (
	"github.com/childrenofukiyo/odin/pkg/httperror"
	"github.com/childrenofukiyo/odin/pkg/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type UserController struct {
	logger      *zap.SugaredLogger
	userService service.UserService
}

func NewUserController(e *echo.Group, logger *zap.SugaredLogger, userService service.UserService) {
	ctrl := &UserController{
		logger:      logger,
		userService: userService,
	}
	e.GET("/hello", ctrl.Hello)
}

func (ctrl *UserController) Hello(c echo.Context) error {
	claims := getClaims(c)

	response, err := ctrl.userService.Get(claims.UserID)
	if err != nil {
		return httperror.FromDomain(err)
	}

	return c.JSON(http.StatusOK, response)
}
