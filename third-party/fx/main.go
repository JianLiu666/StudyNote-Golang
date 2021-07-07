package main

import (
	"context"
	"jian6/third-party/fx/config"
	"jian6/third-party/fx/delivery/server"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/fx"
)

func NewLogger() *log.Logger {
	logger := log.New(os.Stdout, "[FxExample] ", 0)
	logger.Print("Successfully initialized logger.")

	return logger
}

func main() {
	config.NewConfig()

	app := fx.New(
		fx.Supply(
			config.GetConfig(),
		),

		server.Module,

		fx.Provide(
			NewLogger,
		),

		fx.Invoke(
		//
		),
	)

	if err := app.Start(context.TODO()); err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)
	if _, err := http.Get("http://localhost:8001/"); err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)
	if err := app.Stop(context.TODO()); err != nil {
		log.Fatal(err)
	}
}
