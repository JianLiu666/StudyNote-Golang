package server

import (
	"context"
	"fmt"
	"jian6/third_party/fx/config"
	"log"
	"net/http"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		newWithConfig,
	),
	fx.Invoke(
		setRoutes,
		start,
	),
)

type Param struct {
	fx.In
	Config *config.Config
}

func newWithConfig(p Param, logger *log.Logger) *http.Server {
	defer logger.Print("Successfully initialized server.")

	server := &http.Server{
		Addr: fmt.Sprintf(":%v", p.Config.Server.Port),
	}

	return server
}

func setRoutes(s *http.Server, logger *log.Logger) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Print("Got a request.")
	})
}

func start(lc fx.Lifecycle, s *http.Server, logger *log.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			defer logger.Print("started http server.")
			go func() {
				if err := s.ListenAndServe(); err != nil {
					logger.Fatalf("Failed to listen and serve: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			defer logger.Print("showdown http server.")
			if err := s.Shutdown(ctx); err != nil {
				logger.Fatalf("Failed to shudwon server: %v", err)
			}

			return nil
		},
	})
}
