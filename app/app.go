package app

import (
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	onStart func() error
	onStop func() error
}

func NewApp() *App {
	return &App{}
}

func (a *App) OnStart(onStart func() error) {
	a.onStart = onStart
}

func (a *App) OnStop(onStop func() error) {
	a.onStop = onStop
}

func (a *App) Start() {
	signalCh := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<- signalCh
		if err := a.onStop(); err != nil {
			panic(err)
		}
		done <- true
	}()

	go func() {
		if err := a.onStart(); err != nil {
			panic(err)
		}
	}()

	<- done
}
