package api

import (
	"interview20231116/api/router/v1/head"
	"interview20231116/api/router/v1/page"
	"interview20231116/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
)

type server struct {
	app  *fiber.App
	conf config.ServerOpts
}

func Init(conf config.ServerOpts) *server {
	app := fiber.New()

	// set web server logger format
	app.Use(logger.New(logger.Config{
		Format: "[${time}] | ${ip} | ${latency} | ${status} | ${method} | ${path} | Req: ${body} | Resp: ${resBody}\n",
	}))

	api := app.Group("/api")
	head.NewRouter(api)
	page.NewRouter(api)

	return &server{
		app:  app,
		conf: conf,
	}
}

func (s *server) Run() {
	go func() {
		if err := s.app.Listen(s.conf.Port); err != nil {
			logrus.Panicf("starting server on %s failed: %s", s.conf.Port, err.Error())
		}
	}()
}

func (s *server) Shutdown() {
	if err := s.app.Shutdown(); err != nil {
		logrus.Errorf("main: shuting server down failed: %v\n", err.Error())
	}
}
