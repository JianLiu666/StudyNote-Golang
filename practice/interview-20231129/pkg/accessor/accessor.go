package accessor

import (
	"context"
	"interview20231129/pkg/config"
	"sync"

	"github.com/sirupsen/logrus"
)

type shutdownHandler func(context.Context)

type Accessor struct {
	shutdownOnce     sync.Once
	shutdownHandlers []shutdownHandler

	Config *config.Config
}

func Build() *Accessor {
	return &Accessor{
		Config: config.NewFromViper(),
	}
}

func (a *Accessor) Close(ctx context.Context) {
	a.shutdownOnce.Do(func() {
		logrus.Info("start to close accessors.")
		for _, f := range a.shutdownHandlers {
			f(ctx)
		}
	})

	logrus.Info("all accessors closed.")
}
