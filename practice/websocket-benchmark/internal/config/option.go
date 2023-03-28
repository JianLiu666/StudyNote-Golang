package config

type ServerOpts struct {
	Addr string `mapstructure:"addr" yaml:"addr"`
	Port string `mapstructure:"port" yaml:"port"`
}
