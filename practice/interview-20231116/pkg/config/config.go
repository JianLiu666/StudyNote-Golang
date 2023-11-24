package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerOpts `mapstructure:"server"`
	Redis  RedisOpts  `mapstructure:"redis"`
}

func NewFromViper() *Config {
	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("viper.ReadInConfig failed, used default configuration: %v", err)
		return NewFromDefault()
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		logrus.Errorf("viper.Unmarshal failed, used default configuration: %v", err)
		return NewFromDefault()
	}

	return cfg
}

func NewFromDefault() *Config {
	server := ServerOpts{
		Name: "apiserver",
		Port: "6600",
	}

	redis := RedisOpts{
		Address:            "localhost:6379",
		Password:           "",
		DB:                 0,
		PoolSize:           10,
		ListCollectionName: "list",
		PageExpirationTime: 86400,
		PageKeyPrefix:      "page/",
	}

	cfg := &Config{
		Server: server,
		Redis:  redis,
	}

	return cfg
}
