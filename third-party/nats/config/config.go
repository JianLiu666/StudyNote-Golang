package config

import (
	"sync"

	"github.com/spf13/viper"
)

var (
	setOnce sync.Once
	Nats    *config
)

type config struct {
	Addr                   string `mapstructure:"addr" yaml:"addr,omitempty"`
	ReconnInterval         int64  `mapstructure:"reconnInterval" yaml:"reconnInterval,omitempty"`
	ConnectTimeOut         int64  `mapstructure:"connectTimeOut" yaml:"connectTimeOut,omitempty"`
	StanClusterId          string `mapstructure:"stanClusterId" yaml:"stanClusterId,omitempty"`
	StanPingsInterval      int    `mapstructure:"stanPingsInterval" yaml:"stanPingsInterval,omitempty"`
	StanPingsMaxOut        int    `mapstructure:"stanPingsMaxOut" yaml:"stanPingsMaxOut,omitempty"`
	BenchNumProducers      int    `mapstructure:"benchNumProducers" yaml:"benchNumProducers,omitempty"`
	BenchProducerEachTimes int    `mapstructure:"benchProducerEachTimes" yaml:"benchProducerEachTimes,omitempty"`
}

func NewFromViper() (*config, error) {
	var c config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

func SetConfig(c *config) {
	setOnce.Do(func() {
		Nats = c
	})
}
