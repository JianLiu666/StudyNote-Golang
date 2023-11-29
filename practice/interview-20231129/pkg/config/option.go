package config

type ServerOpts struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
}
