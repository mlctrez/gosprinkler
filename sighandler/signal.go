package sighandler

import (
	"os"
	"os/signal"
	"syscall"
)

type Api struct {
	signalChan chan os.Signal
}

func New() *Api {
	return &Api{signalChan: make(chan os.Signal)}
}

func (a *Api) RegisterHandler(callback func()) {

	signal.Notify(a.signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		<-a.signalChan
		callback()
		os.Exit(1)
	}()

}

func (a *Api) Interrupt() {
	a.signalChan <- os.Interrupt
}
