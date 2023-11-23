package api

import (
	"interview20231116/api/router/v1/list"
	"interview20231116/api/router/v1/page"
	"interview20231116/pkg/accessor"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
)

type server struct {
	app   *fiber.App
	infra *accessor.Accessor
}

func Init(infra *accessor.Accessor) *server {
	app := fiber.New()

	// set web server logger format
	app.Use(logger.New(logger.Config{
		Format: "[${time}] | ${ip} | ${latency} | ${status} | ${method} | ${path} | Req: ${body} | Resp: ${resBody}\n",
	}))

	api := app.Group("/api")
	list.NewListRouter(infra.KvStore).Init(api)
	page.NewPageRouter(infra.KvStore).Init(api)

	return &server{
		app:   app,
		infra: infra,
	}
}

func (s *server) Run() {
	go func() {
		if err := s.app.Listen(s.infra.Config.Server.Port); err != nil {
			logrus.Panicf("starting server on %s failed: %s", s.infra.Config.Server.Port, err.Error())
		}
	}()
}

func (s *server) Shutdown() {
	if err := s.app.Shutdown(); err != nil {
		logrus.Errorf("main: shuting server down failed: %v\n", err.Error())
	}
}
