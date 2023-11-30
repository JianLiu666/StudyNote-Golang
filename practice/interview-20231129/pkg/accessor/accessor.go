package accessor

import (
	"context"
	"interview20231129/pkg/config"
	"interview20231129/pkg/singlepool"
	"sync"

	"github.com/sirupsen/logrus"
)

type shutdownHandler func(context.Context)

type Accessor struct {
	shutdownOnce     sync.Once
	shutdownHandlers []shutdownHandler

	Config     *config.Config
	SinglePool singlepool.SinglePool
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

func (a *Accessor) InitSinglePool(ctx context.Context, singlePool singlepool.SinglePool) {
	a.SinglePool = singlePool

	a.shutdownHandlers = append(a.shutdownHandlers, func(c context.Context) {
		// TODO: graceful shutdown
		logrus.Infoln("single pool accessor closed.")
	})

	logrus.Infoln("initial single pool accessor successful.")
}
