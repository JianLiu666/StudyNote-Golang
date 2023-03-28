package config

import "github.com/spf13/viper"

var cfg *Config

type Config struct {
	Server ServerOpts `mapstructure:"server" yaml:"server"`
}

func NewFromViper() *Config {
	err := viper.ReadInConfig()
	if err != nil {
		return NewFromDefault()
	}

	cfg = &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return NewFromDefault()
	}

	return cfg
}

func NewFromDefault() *Config {
	server := ServerOpts{
		Addr: "localhost",
		Port: "6600",
	}

	cfg := &Config{
		Server: server,
	}

	return cfg
}
