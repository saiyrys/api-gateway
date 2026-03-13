package app

import (
	"api-gateway/api"
	"api-gateway/pkg/exceptions"
	"api-gateway/pkg/graceful"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

type App struct {
	mainCtx    context.Context
	router     *http.ServeMux
	httpServer *http.Server
	signalChan chan os.Signal
}

func (a *App) AddHandler(path string, handler http.HandlerFunc) *App {
	if a.router != nil {
		a.router.HandleFunc(path, handler)
	}
	return a
}

func NewApp(mainCtx context.Context) *App {
	return &App{
		mainCtx:    mainCtx,
		signalChan: make(chan os.Signal, 1),
		router: 		http.NewServeMux(),
	}
}

func (a *App) WithHTTPServer(addr string) *App {
	a.httpServer = &http.Server{
		Addr:         addr,
		Handler:      a.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return a
}

func (a *App) InitServer(addr string) error {
	sb := api.NewServerBuilder(a.mainCtx)

	sb.AddHandlers()

	a.router = sb.GetRouter()

	return nil
}

func (a *App) Run(se *exceptions.ServerErrors) error {
	if a.httpServer == nil {
		return errors.New("HTTP server not initialized")
	}

	shutdown := graceful.NewGracefulShutdown()

	go func() {
		log.Printf("HTTP server starting on %s", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {

			se.AddServerErrors(err)
		}
	}()

	shutdown.WaitForShutdown(se)
	log.Println("Получен сигнал завершения, приложение останавливается")
	return a.Stop(se)
}

func (a *App) Stop(se *exceptions.ServerErrors) error {
	if a.httpServer == nil {
		return nil
	}

	log.Println("Сервер завершает работу...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		se.AddServerErrors(err)
		return err
	}

	allError := se.GetErrors()
	for i, err := range allError {
		log.Printf("Ошибка %d: %v", i+1, err)
	}

	log.Println("Сервер Остановлен")
	return nil
}
