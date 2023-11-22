package config

type ServerOpts struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
}

type RedisOpts struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}
