package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerOpts `mapstructure:"server"`
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

	cfg := &Config{
		Server: server,
	}

	return cfg
}
