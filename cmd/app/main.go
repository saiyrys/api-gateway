package main

import (
	"api-gateway/internal/app"
	"api-gateway/pkg/exceptions"
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	errors := exceptions.NewServerErrors(10)
	defer errors.Close()

	application := app.NewApp(ctx)
	
	if err := application.InitServer(":8080"); err != nil {
		log.Fatal("Ошибка инициализации сервера:", err)
	}

	application.WithHTTPServer(":8080")
	
	if err := application.Run(errors); err != nil {
		log.Fatal("Ошибка запуска:", err)
	}
}