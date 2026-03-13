package graceful

import (
	"api-gateway/pkg/exceptions"
	"os"
	"os/signal"
	"syscall"
)

type GracefulShutdown struct {
	signalChan    			chan os.Signal
}

func NewGracefulShutdown() *GracefulShutdown{
	return &GracefulShutdown{
		signalChan: make(chan os.Signal, 1),
	}
}

// WaitForShutdown Метод ожидает сигнал завершения и добавляет ошибки
func (g *GracefulShutdown) WaitForShutdown(se *exceptions.ServerErrors) {
	// Региструем сигнал для завершения
	signal.Notify(g.signalChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Блокируемся до получения сигнала
	<-g.signalChan

	// Добавляем ошибку о завершении работы
	se.AddServerErrors(os.ErrProcessDone)
}
