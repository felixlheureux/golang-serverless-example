package controller

import (
	"github.com/childrenofukiyo/odin/pkg/auth"
	"github.com/childrenofukiyo/odin/pkg/httperror"
	"github.com/childrenofukiyo/odin/pkg/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type AuthController struct {
	logger      *zap.SugaredLogger
	authService service.AuthService
}

func NewAuthController(e *echo.Group, logger *zap.SugaredLogger, authService service.AuthService, authenticator echo.MiddlewareFunc) {
	ctrl := &AuthController{
		logger:      logger,
		authService: authService,
	}
	e.POST("/challenge", ctrl.Challenge)
	e.POST("/authorize", ctrl.Authorize)
	e.POST("/authorize/silently", ctrl.AuthorizeSilently, authenticator)
}

func (ctrl *AuthController) Challenge(c echo.Context) error {
	addressHex := c.FormValue("ethereum_address")

	input := auth.NewChallengeInput(addressHex)

	response, err := ctrl.authService.Challenge(input)
	if err != nil {
		return httperror.FromDomain(err)
	}

	return c.JSON(http.StatusOK, response)
}

func (ctrl *AuthController) Authorize(c echo.Context) error {
	addressHex := c.FormValue("ethereum_address")
	sigHex := c.FormValue("signature")

	input := auth.NewAuthorizeInput(addressHex, sigHex)

	response, err := ctrl.authService.Authorize(input)
	if err != nil {
		return httperror.FromDomain(err)
	}

	return c.JSON(http.StatusOK, response)
}

func (ctrl *AuthController) AuthorizeSilently(c echo.Context) error {
	claims := getClaims(c)

	input := auth.NewAuthorizeSilentlyInput(claims.EthereumAddressHex)

	response, err := ctrl.authService.AuthorizeSilently(input)
	if err != nil {
		return httperror.FromDomain(err)
	}

	return c.JSON(http.StatusOK, response)
}
