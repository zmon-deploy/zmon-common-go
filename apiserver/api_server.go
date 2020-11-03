package apiserver

import (
	"context"
	"fmt"
	"github.com/cnpst/zmon-common-go/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type Controller interface {
	Route(e *echo.Echo)
}

func StartApiServer(logger log.Logger, controllers []Controller, port int) (func(), error) {
	e := echo.New()
	e.HTTPErrorHandler = getErrorHandler(logger)
	e.Use(middleware.Recover())
	e.Logger.SetOutput(ioutil.Discard)

	for _, controller := range controllers {
		controller.Route(e)
	}

	if err := e.Start(fmt.Sprintf(":%d", port)); err != nil {
		return nil, errors.Wrap(err, "Failed to start echo server")
	}

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			logger.Errorf("failed to shutdown api server: %s", err.Error())
		}
	}, nil
}

func getErrorHandler(logger log.Logger) func(err error, ctx echo.Context) {
	return func(err error, ctx echo.Context) {
		code := http.StatusInternalServerError
		if httpError, ok := err.(*echo.HTTPError); ok {
			code = httpError.Code
		}

		logger.Errorf("error occurred, code: %d, err: %s", code, err.Error())

		_ = ctx.JSON(code, map[string]interface{}{
			"ok":      false,
			"message": err.Error(),
		})
	}
}
