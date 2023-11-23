package accessor

import (
	"context"
	"interview20231116/pkg/config"
	"interview20231116/pkg/kvstore"
	"sync"

	"github.com/sirupsen/logrus"
)

type shutdownHandler func(context.Context)

type Accessor struct {
	shutdownOnce     sync.Once
	shutdownHandlers []shutdownHandler

	Config  *config.Config
	KvStore kvstore.KvStore
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

func (a *Accessor) InitKvStore(ctx context.Context) {
	a.KvStore = kvstore.NewRedisClient(ctx, &a.Config.Redis)

	a.shutdownHandlers = append(a.shutdownHandlers, func(c context.Context) {
		a.KvStore.Shutdown(c)
		logrus.Infoln("key-value store accessor closed.")
	})

	logrus.Infoln("initial key-value store accessor successful.")
}
