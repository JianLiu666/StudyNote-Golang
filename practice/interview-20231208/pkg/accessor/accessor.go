package accessor

import (
	"context"
	"fmt"
	"interview20231208/pkg/config"
	"interview20231208/pkg/rdb"
	"interview20231208/pkg/rdb/mysql"
	"interview20231208/pkg/trading"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type shutdownHandler func(context.Context)

type Accessor struct {
	shutdownOnce     sync.Once
	shutdownHandlers []shutdownHandler

	Config      *config.Config
	RDB         rdb.RDB
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

func (a *Accessor) InitRDB(ctx context.Context) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		a.Config.MySQL.UserName,
		a.Config.MySQL.Password,
		a.Config.MySQL.Address,
		a.Config.MySQL.DBName,
	)

	a.RDB = mysql.NewMySqlClient(ctx, dsn,
		time.Duration(a.Config.MySQL.ConnMaxLifetime)*time.Minute,
		a.Config.MySQL.MaxOpenConns,
		a.Config.MySQL.MaxIdleConns,
	)

	a.shutdownHandlers = append(a.shutdownHandlers, func(c context.Context) {
		a.RDB.Shutdown(c)
		logrus.Infoln("relational database accessor closed.")
	})

	logrus.Infoln("initial relational database accessor successful.")
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
