package httperror

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Response struct {
	ErrorCode int    `json:"error_code"`
	Error     string `json:"error"`
}

type Error struct {
	StatusCode int
	ErrorCode  int
	// message returned to the client
	OutputMessage string
	Cause         error
}

func (e Error) Error() string {
	return fmt.Sprintf("[%d] %s: %+v", e.ErrorCode, e.OutputMessage, e.Cause)
}

func (e Error) Is(other error) bool {
	err, ok := other.(*Error)

	if !ok {
		return false
	}

	return err.ErrorCode == e.ErrorCode
}

func NewErrorHandler(logger *zap.SugaredLogger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var e *Error

		switch err.(type) {
		case *echo.HTTPError:
			e = CoreEchoError(err.(*echo.HTTPError))
		case *Error:
			e = err.(*Error)
		default:
			e = CoreUnknownError(err)
		}

		if c.Response().Committed {
			logger.Warnw("response already committed", zap.Error(err))
			return
		}

		if c.Request().Method == http.MethodHead { // Issue https://github.com/labstack/echo/issues/608
			err = c.NoContent(e.StatusCode)
		} else {
			err = c.JSON(e.StatusCode, Response{e.ErrorCode, e.OutputMessage})
		}

		if err != nil {
			logger.Errorw("unable to write error response", zap.Error(err))
		}
	}
}
