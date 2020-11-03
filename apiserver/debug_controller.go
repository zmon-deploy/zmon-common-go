package apiserver

import (
	"github.com/labstack/echo/v4"
	"net/http/pprof"
)

type DebugController struct {
}

func (c *DebugController) Route(e *echo.Echo) {
	e.GET("/debug/pprof", c.index)
	e.GET("/debug/pprof/heap", c.heap)
	e.GET("/debug/pprof/goroutine", c.goroutine)
	e.GET("/debug/pprof/block", c.block)
	e.GET("/debug/pprof/threadcreate", c.threadCreate)
	e.GET("/debug/pprof/cmdline", c.cmdline)
	e.GET("/debug/pprof/profile", c.profile)
	e.GET("/debug/pprof/symbol", c.symbol)
	e.POST("/debug/pprof/symbol", c.symbol)
	e.GET("/debug/pprof/trace", c.trace)
	e.GET("/debug/pprof/mutex", c.mutex)
}

func (c *DebugController) index(ctx echo.Context) error {
	pprof.Index(ctx.Response().Writer, ctx.Request())
	return nil
}

func (c *DebugController) heap(ctx echo.Context) error {
	pprof.Handler("heap").ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

func (c *DebugController) goroutine(ctx echo.Context) error {
	pprof.Handler("goroutine").ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

func (c *DebugController) block(ctx echo.Context) error {
	pprof.Handler("block").ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

func (c *DebugController) threadCreate(ctx echo.Context) error {
	pprof.Handler("threadcreate").ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

func (c *DebugController) cmdline(ctx echo.Context) error {
	pprof.Cmdline(ctx.Response().Writer, ctx.Request())
	return nil
}

func (c *DebugController) profile(ctx echo.Context) error {
	pprof.Profile(ctx.Response().Writer, ctx.Request())
	return nil
}

func (c *DebugController) symbol(ctx echo.Context) error {
	pprof.Symbol(ctx.Response().Writer, ctx.Request())
	return nil
}

func (c *DebugController) trace(ctx echo.Context) error {
	pprof.Trace(ctx.Response().Writer, ctx.Request())
	return nil
}

func (c *DebugController) mutex(ctx echo.Context) error {
	pprof.Handler("mutex").ServeHTTP(ctx.Response().Writer, ctx.Request())
	return nil
}

