package api

import (
	"context"
	"interview20231129/api/router/v1/single"
	"interview20231129/pkg/accessor"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type server struct {
	app   *http.Server
	infra *accessor.Accessor
}

func Init(infra *accessor.Accessor) *server {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	api := router.Group("/api")
	single.NewSingleRouter().Init(api)

	app := &http.Server{
		Addr:    infra.Config.Server.Port,
		Handler: router,
	}

	return &server{
		app:   app,
		infra: infra,
	}
}

func (s *server) Run() {
	go func() {
		if err := s.app.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Panicf("starting server on %s failed: %s", s.infra.Config.Server.Port, err.Error())
		}
	}()
}

func (s *server) Shutdown(ctx context.Context) {
	if err := s.app.Shutdown(ctx); err != nil {
		logrus.Errorf("shuting server down failed: %v\n", err.Error())
	}
}
