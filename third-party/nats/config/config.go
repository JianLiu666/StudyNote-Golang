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
	Addr              string `mapstructure:"addr" yaml:"addr,omitempty"`
	ClusterId         string `mapstructure:"clusterId" yaml:"clusterId,omitempty"`
	ReconnInterval    int64  `mapstructure:"reconnInterval" yaml:"reconnInterval,omitempty"`
	ConnectTimeOut    int64  `mapstructure:"connectTimeOut" yaml:"connectTimeOut,omitempty"`
	StanPingsInterval int    `mapstructure:"stanPingsInterval" yaml:"stanPingsInterval,omitempty"`
	StanPingsMaxOut   int    `mapstructure:"stanPingsMaxOut" yaml:"stanPingsMaxOut,omitempty"`
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
