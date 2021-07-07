package config

import "sync"

var initialized sync.Once
var globalConfig *Config

type Config struct {
	Server *ServerConf
}

func NewConfig() {
	initialized.Do(func() {
		serverConf := &ServerConf{
			Port: "8001",
		}

		globalConfig = &Config{
			Server: serverConf,
		}
	})
}

func GetConfig() *Config {
	return globalConfig
}

type ServerConf struct {
	Port string
}
