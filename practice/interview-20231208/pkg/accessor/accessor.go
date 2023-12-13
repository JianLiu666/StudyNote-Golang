package accessor

import (
	"context"
	"interview20231208/pkg/config"
	"interview20231208/pkg/trading"
	"sync"

	"github.com/sirupsen/logrus"
)

type shutdownHandler func(context.Context)

type Accessor struct {
	shutdownOnce     sync.Once
	shutdownHandlers []shutdownHandler

	Config      *config.Config
	TradingPool trading.TradingPool
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

func (a *Accessor) InitTradingPool(ctx context.Context, tradingPool trading.TradingPool) {
	a.TradingPool = tradingPool

	a.shutdownHandlers = append(a.shutdownHandlers, func(c context.Context) {
		// TODO: graceful shutdown
		logrus.Infoln("trading pool accessor closed.")
	})

	// TODO: consider how to control context
	a.TradingPool.Enable(ctx)

	logrus.Infoln("initial trading pool accessor successful.")
}
