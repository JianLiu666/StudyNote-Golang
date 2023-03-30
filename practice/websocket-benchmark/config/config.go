package config

import "github.com/spf13/viper"

var cfg *Config

type Config struct {
	Server     ServerOpts     `yaml:"server"`
	Simulation SimulationOpts `yaml:"simulation"`
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

	simulation := SimulationOpts{
		NumClients:  1,
		NumMessages: 100,
	}

	cfg := &Config{
		Server:     server,
		Simulation: simulation,
	}

	return cfg
}
