package accessor

import (
	"context"
	"interview20231116/pkg/config"
	"sync"

	"github.com/sirupsen/logrus"
)

type shutdownHandler func(context.Context)

type accessor struct {
	shutdownOnce     sync.Once
	shutdownHandlers []shutdownHandler

	Config *config.Config
}

func Build() *accessor {
	return &accessor{
		Config: config.NewFromViper(),
	}
}

func (a *accessor) Close(ctx context.Context) {
	a.shutdownOnce.Do(func() {
		logrus.Info("start to close accessors.")
		for _, f := range a.shutdownHandlers {
			f(ctx)
		}
	})

	logrus.Info("all accessors closed.")
}
